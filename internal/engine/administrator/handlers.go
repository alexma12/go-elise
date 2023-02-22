package administrator

import (
	"fmt"
	"net/http"

	"github.com/alexma12/go-elise/internal/model"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func (a *Admin) addConfig(c echo.Context) error {
	var config model.ScrapeConfig
	if err := c.Bind(&config); err != nil {
		return c.String(http.StatusBadRequest, "Invalid scrapeConfig")
	}
	id := uuid.New()
	err := a.DB.AddConfig(id, config.Name, config.Url, config.Selector, config.Type, config.RequiresWebDriver)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Could not create config")
	}
	config.ID = id
	a.InfoLog.Printf("created config with id : %s", id)
	return c.JSON(http.StatusCreated, config)
}

func (a *Admin) listConfigs(c echo.Context) error {
	configs, err := a.DB.ListConfigs()
	if err != nil {
		return c.String(http.StatusInternalServerError, "Could not get configs")
	}

	a.InfoLog.Println(configs)
	return c.JSON(http.StatusOK, configs)
}

func (a *Admin) deleteConfig(c echo.Context) error {
	idParam := c.QueryParams().Get("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		a.ErrorLog.Println(err)
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Invalid ID received %s", idParam))
	}
	fmt.Println(id)

	deleted, err := a.DB.DeleteConfig(id)
	if err != nil {
		a.ErrorLog.Println(err)
		return c.String(http.StatusInternalServerError, "Could not delete config")
	}

	if deleted {
		a.InfoLog.Printf("Admin: Deleted Config with ID: %s", id)
		return c.JSON(http.StatusOK, "Successfully deleted")
	} else {
		return c.String(http.StatusNotFound, fmt.Sprintf("Config with ID: %s does not exist", id))
	}
}