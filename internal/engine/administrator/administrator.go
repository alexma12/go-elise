// package admin contains the admin module code
// The purpose of the admin module is to manage and relay information across modules
package administrator

import (
	"log"

	"github.com/alexma12/go-elise/internal/model"
	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
)

type Admin struct {
	DB       model.ScrapeStorage
	ErrorLog *log.Logger
	InfoLog  *log.Logger
}

func New(db model.ScrapeStorage, errorLog, infoLog *log.Logger) *Admin {
	return &Admin{
		DB:       db,
		ErrorLog: errorLog,
		InfoLog:  infoLog,
	}
}

func (a *Admin) Start() {
	err := a.DB.CreateStorage()
	if err != nil {
		a.ErrorLog.Println(err)
		return
	}
	server := a.initServer()
	a.InfoLog.Printf("Admin: Starting server on port: %s", "3030")
	a.ErrorLog.Fatal(server.Start(":3030"))
}

func (a *Admin) initServer() *echo.Echo {
	e := echo.New()

	g := e.Group("/elise/config")
	g.POST("/add", a.addConfig)
	g.GET("/list", a.listConfigs)
	g.DELETE("/delete", a.deleteConfig)

	return e
}
