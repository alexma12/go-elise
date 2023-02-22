package scrapedb

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
	RequiresWebDriver bool
	CreatedAt         time.Time
}

type ScrapeLog struct {
	ID         uuid.UUID
	Value      string
	ExecutedAt time.Time
}

type ScrapeDB interface {
	CreateTables() error
	AddConfig(id uuid.UUID, name, url, selector string, targetType TargetType, requiresWebDriver bool) error
	ListConfigs() ([]ScrapeConfig, error)
	DeleteConfig(id uuid.UUID) (bool, error)
}
