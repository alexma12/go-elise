package scheduler

import (
	"context"

	"github.com/alexma12/go-elise/pkg/model"
)

type jobExecutor struct{}

func NewJobExecutor() *jobExecutor {
	return &jobExecutor{}
}

func (j *jobExecutor) Start(ctx context.Context) chan<- model.ScrapeConfig {
	toExecute := make(chan model.ScrapeConfig)
	go func() {
		defer close(toExecute)
	}()
	return toExecute
}
