package mysql

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/alexma12/go-elise/pkg/model"
	"github.com/google/uuid"
)

type MySQLScrapeConfig struct {
	ID                uuid.UUID
	Name              string
	Url               string
	Selector          string
	Type              int
	Interval          int
	RequiresWebDriver bool
	CreatedAt         time.Time
}

func (sc MySQLScrapeConfig) toScrapeConfig() model.ScrapeConfig {
	return model.ScrapeConfig{
		ID:        sc.ID,
		Name:      sc.Name,
		Url:       sc.Url,
		Selector:  sc.Selector,
		Type:      model.TargetType(sc.Type),
		Interval:  sc.Interval,
		CreatedAt: sc.CreatedAt,
	}
}

type MySQLScrapeLog struct {
	ID         uuid.UUID
	Value      string
	ExecutedAt time.Time
}

func (sl MySQLScrapeLog) toScrapeLog() model.ScrapeLog {
	return model.ScrapeLog{
		ID:         sl.ID,
		Value:      sl.Value,
		ExecutedAt: sl.ExecutedAt,
	}
}

type MySQLScrapeDB struct {
	DB *sql.DB
}

func New(db *sql.DB) *MySQLScrapeDB {
	return &MySQLScrapeDB{DB: db}
}

func (ms *MySQLScrapeDB) CreateTables() error {
	stmt := `USE elise;`
	_, err := ms.DB.Exec(stmt)
	if err != nil {
		return err
	}

	fmt.Println("creating scraper scrape config table, if not exists")
	stmt = `CREATE TABLE IF NOT EXISTS scrape_configs (
        id BINARY(16) NOT NULL PRIMARY KEY,
        name VARCHAR(255) NOT NULL, 
        url VARCHAR(2048) NOT NULL,
        selector VARCHAR(1024) NOT NULL,
        type INT NOT NULL,
		interval_minutes INT NOT NULL,
        requiresWebDriver BOOLEAN NOT NULL,
        createdAt TIMESTAMP NOT NULL
    )`
	_, err = ms.DB.Exec(stmt)
	if err != nil {
		return err
	}

	fmt.Println("creating scraper scrape logs table, if not exists")
	stmt = `CREATE TABLE IF NOT EXISTS scrape_logs (
		id BINARY(16) NOT NULL PRIMARY KEY,
		value VARCHAR(2048), 
		executedAt TIMESTAMP NOT NULL
	)`
	_, err = ms.DB.Exec(stmt)
	if err != nil {
		return err
	}

	return nil
}

func (ms *MySQLScrapeDB) AddConfig(id uuid.UUID, name, url, selector string, targetType model.TargetType, interval int, requiresWebDriver bool) error {
	stmt := `INSERT INTO scrape_configs (id, name, url, selector, type, interval_minutes, requiresWebDriver, createdAt)
             VALUES(UUID_TO_BIN(?), ?, ?, ?, ?, ?, ?, UTC_TIMESTAMP())`
	_, err := ms.DB.Exec(stmt, id, name, url, selector, targetType, interval, requiresWebDriver)
	if err != nil {
		return err
	}
	return nil
}

func (ms *MySQLScrapeDB) ListConfigs() ([]model.ScrapeConfig, error) {
	stmt := `SELECT * FROM scrape_configs`
	rows, err := ms.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	configs := []model.ScrapeConfig{}
	for rows.Next() {
		var c MySQLScrapeConfig
		if err := rows.Scan(&c.ID, &c.Name, &c.Url, &c.Selector, &c.Type, &c.Interval, &c.RequiresWebDriver, &c.CreatedAt); err != nil {
			return nil, err
		}
		configs = append(configs, c.toScrapeConfig())
	}

	return configs, nil
}

func (ms *MySQLScrapeDB) DeleteConfig(id uuid.UUID) (bool, error) {
	stmt := `DELETE FROM scrape_configs WHERE id = UUID_TO_BIN(?)`
	res, err := ms.DB.Exec(stmt, id)
	if err != nil {
		return false, err
	}
	i, _ := res.RowsAffected()
	return i >= 1, nil
}
