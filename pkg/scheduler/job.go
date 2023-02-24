package scheduler

import (
	"time"

	"github.com/alexma12/go-elise/pkg/model"
)

func NewScrapeJob(config model.ScrapeConfig) *scrapeJob {
	return &scrapeJob{
		scrapeConfig: config,
		// TODO: currently time.Second for testing purposes, need to change to time.Minute
		executeTime: time.Now().Add(time.Duration(config.Interval) * 10 * time.Second),
	}
}

type scrapeJob struct {
	scrapeConfig model.ScrapeConfig
	executeTime  time.Time
}

func (j *scrapeJob) isDue() bool {
	return j.executeTime.Before(time.Now())
}
