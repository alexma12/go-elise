package scheduler

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/alexma12/go-elise/pkg/model"
	"github.com/google/uuid"
)

type ScrapeScheduler struct {
	mu            sync.Mutex
	scrapeConfigs map[uuid.UUID]model.ScrapeConfig
	// may contain jobs that have been removed
	jobQueue    *jobQueue
	jobExecutor *jobExecutor

	// TODO: rethink about these job channels
	addJobChan     chan<- *scrapeJob
	executeJobChan chan<- *scrapeJob

	errorLog *log.Logger
	infoLog  *log.Logger
}

func New(errorLog, infoLog *log.Logger, scrapeConfigs []model.ScrapeConfig) *ScrapeScheduler {
	configs := map[uuid.UUID]model.ScrapeConfig{}
	for _, c := range scrapeConfigs {
		configs[c.ID] = c
	}

	jobs := make([]*scrapeJob, len(scrapeConfigs))
	for i, c := range scrapeConfigs {
		// TODO: add some sort of randomized delay
		jobs[i] = NewScrapeJob(c)
	}

	addJobChan := make(chan *scrapeJob, 10)
	jobQueue := NewJobQueue(addJobChan, jobs)

	executeJobChan := make(chan *scrapeJob, 10)
	jobExecutor := NewJobExecutor(executeJobChan)

	return &ScrapeScheduler{
		scrapeConfigs:  configs,
		jobQueue:       jobQueue,
		jobExecutor:    jobExecutor,
		addJobChan:     addJobChan,
		executeJobChan: executeJobChan,
		errorLog:       errorLog,
		infoLog:        infoLog,
	}
}

func (s *ScrapeScheduler) Start(ctx context.Context) {
	go func() {
		defer close(s.addJobChan)
		defer close(s.executeJobChan)

		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		readyJobs := s.jobQueue.Start(ctx)
		finishedJobs := s.jobExecutor.Start(ctx)
		for {
			select {
			case <-ctx.Done():
				log.Println("context cancelled")
				return
			case j, ok := <-readyJobs:
				if !ok {
					fmt.Println("ready channel cancelled")
					return
				}
				s.mu.Lock()
				_, ok = s.scrapeConfigs[j.scrapeConfig.ID]
				s.mu.Unlock()

				if ok {
					s.executeJobChan <- j
				}

			case j, ok := <-finishedJobs:
				if !ok {
					fmt.Println("finished channel cancelled")
					return
				}
				s.mu.Lock()
				_, ok = s.scrapeConfigs[j.scrapeConfig.ID]
				s.mu.Unlock()

				if ok {
					newJob := NewScrapeJob(j.scrapeConfig)
					s.addJobChan <- newJob
				}
			}
		}
	}()
}

func (s *ScrapeScheduler) AddJob(scrapeConfig model.ScrapeConfig) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.scrapeConfigs[scrapeConfig.ID]; ok {
		return
	}
	s.scrapeConfigs[scrapeConfig.ID] = scrapeConfig
	s.addJobChan <- NewScrapeJob(scrapeConfig)
}

func (s *ScrapeScheduler) GetCurrentJobs() []scrapeJob {
	s.mu.Lock()
	defer s.mu.Unlock()

	out := []scrapeJob{}
	for _, j := range *s.jobQueue.pending {
		// TODO add some status to indicate processing/pending ?
		if _, ok := s.scrapeConfigs[j.scrapeConfig.ID]; ok {
			out = append(out, *j)
		}
	}
	return out
}

func (s *ScrapeScheduler) DeleteJob(id uuid.UUID) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.scrapeConfigs[id]; !ok {
		return
	}
	delete(s.scrapeConfigs, id)
}
