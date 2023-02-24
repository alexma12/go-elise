package scheduler

import (
	"context"
	"fmt"
	"time"
)

type jobExecutor struct {
	executeJobChan <-chan *scrapeJob
}

func NewJobExecutor(executeJobChan <-chan *scrapeJob) *jobExecutor {
	return &jobExecutor{
		executeJobChan: executeJobChan,
	}
}

func (j *jobExecutor) Start(ctx context.Context) <-chan *scrapeJob {
	finished := make(chan *scrapeJob, 10)
	go func() {
		fmt.Println("START executed")
		defer close(finished)
		for {
			select {
			case <-ctx.Done():
				fmt.Println("context cancelled")
				return
			case exec, ok := <-j.executeJobChan:
				if !ok {
					fmt.Println("execute job channel cancelled")
				}
				time.Sleep(8 * time.Second)
				fmt.Printf("Finished executing: %v \n", exec.scrapeConfig.ID)
				finished <- exec
			}
		}
	}()
	return finished
}
