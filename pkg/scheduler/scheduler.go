package scheduler

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/alexma12/go-elise/pkg/model"
	"github.com/google/uuid"
)

type scrapeScheduler struct {
	mu            sync.Mutex
	scrapeConfigs map[uuid.UUID]model.ScrapeConfig
	// may contain jobs that have been removed
	jobQueue    *jobQueue
	jobExecutor *jobExecutor
	addJobChan  chan<- uuid.UUID
	errorLog    *log.Logger
	infoLog     *log.Logger
}

func New(errorLog, infoLog *log.Logger, scrapeConfigs []model.ScrapeConfig) *scrapeScheduler {
	configs := map[uuid.UUID]model.ScrapeConfig{}
	for _, c := range scrapeConfigs {
		configs[c.ID] = c
	}

	jobs := make([]scrapeJob, len(scrapeConfigs))
	for i, c := range scrapeConfigs {
		// TODO: add some sort of randomized delay
		jobs[i] = scrapeJob{
			id:          c.ID,
			executeTime: time.Now().Add(10 * time.Second),
		}
	}

	addJobChan := make(chan uuid.UUID)
	jobQueue := NewJobQueue(addJobChan, jobs)
	// jobExecutor := NewJobExecutor()

	return &scrapeScheduler{
		scrapeConfigs: configs,
		jobQueue:      jobQueue,
		// jobExecutor:   jobExecutor,
		addJobChan: addJobChan,
		errorLog:   errorLog,
		infoLog:    infoLog,
	}
}

func (s *scrapeScheduler) Start(ctx context.Context) {
	go func() {
		defer close(s.addJobChan)
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		readyJobs := s.jobQueue.Start(ctx)
		for {
			select {
			case <-ctx.Done():
				log.Println("cancelled")
			case j := <-readyJobs:
				s.mu.Lock()
				if c, ok := s.scrapeConfigs[j]; ok {
					fmt.Println(c)
				}
				s.mu.Unlock()
			}
		}
	}()
}

func (s *scrapeScheduler) AddJob(scrapeConfig model.ScrapeConfig) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.scrapeConfigs[scrapeConfig.ID]; ok {
		return
	}
	s.scrapeConfigs[scrapeConfig.ID] = scrapeConfig
	s.addJobChan <- scrapeConfig.ID
}

func (s *scrapeScheduler) GetCurrentJobs() []scrapeJob {
	s.mu.Lock()
	defer s.mu.Unlock()

	out := []scrapeJob{}
	for _, j := range *s.jobQueue.pending {
		// TODO add some status to indicate processing/pending ?
		if _, ok := s.scrapeConfigs[j.id]; ok {
			out = append(out, j)
		}
	}
	return out
}

func (s *scrapeScheduler) DeleteJob(id uuid.UUID) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.scrapeConfigs[id]; !ok {
		return
	}
	delete(s.scrapeConfigs, id)
}
