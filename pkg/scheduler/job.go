package scheduler

import (
	"time"

	"github.com/google/uuid"
)

type scrapeJob struct {
	id          uuid.UUID
	executeTime time.Time
}

func (j scrapeJob) isDue() bool {
	return j.executeTime.Before(time.Now())
}
