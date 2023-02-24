// package admin contains the admin module code
// The purpose of the admin module is to manage and relay information across modules
package administrator

import (
	"context"
	"log"

	"github.com/alexma12/go-elise/pkg/model"
	"github.com/alexma12/go-elise/pkg/scheduler"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type scrapeDB interface {
	CreateTables() error
	AddConfig(id uuid.UUID, name, url, selector string, targetType model.TargetType, interval int, requiresWebDriver bool) error
	ListConfigs() ([]model.ScrapeConfig, error)
	DeleteConfig(id uuid.UUID) (bool, error)
}

type admin struct {
	scrapeDB        scrapeDB
	scrapeScheduler *scheduler.ScrapeScheduler
	errorLog        *log.Logger
	infoLog         *log.Logger
}

func New(db scrapeDB, errorLog, infoLog *log.Logger) *admin {
	return &admin{
		scrapeDB:        db,
		scrapeScheduler: scheduler.New(errorLog, infoLog),
		errorLog:        errorLog,
		infoLog:         infoLog,
	}
}

func (a *admin) Start() {
	err := a.scrapeDB.CreateTables()
	if err != nil {
		a.errorLog.Println(err)
		return
	}

	a.infoLog.Println("Admin: Initializing Scheduler...")
	configs, err := a.scrapeDB.ListConfigs()
	if err != nil {
		a.errorLog.Printf("error occured when fetching all configs %s", err)
		configs = []model.ScrapeConfig{}
	}
	a.scrapeScheduler.Start(context.TODO(), configs)

	server := a.initServer()
	a.infoLog.Printf("Admin: Starting server on port: %s", "3030")
	a.errorLog.Fatal(server.Start(":3030"))
}

func (a *admin) initServer() *echo.Echo {
	e := echo.New()

	g := e.Group("/elise/config")
	g.POST("/add", a.addConfig)
	g.GET("/list", a.listConfigs)
	g.DELETE("/delete", a.deleteConfig)

	return e
}
