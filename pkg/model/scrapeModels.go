package model

import (
	"time"

	"github.com/google/uuid"
)

type TargetType int

const (
	Number TargetType = iota
	String
	// TODO: add more
)

type ScrapeConfig struct {
	ID                uuid.UUID
	Name              string
	Url               string
	Selector          string
	Type              TargetType
	Interval          int
	RequiresWebDriver bool
	CreatedAt         time.Time
}

type ScrapeLog struct {
	ID         uuid.UUID
	Value      string
	ExecutedAt time.Time
}
