package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time-progression/progress"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type settings struct {
	HostAddress     string
	HostPort        int
	DefaultTimezone string
	DefaultFormat   string
}

var config settings

func init() {
	file, err := os.Open("config.json")
	if err != nil {
		fmt.Println("Error opening config file:", err)
		os.Exit(1)
	}

	decoder := json.NewDecoder(file)
	if err = decoder.Decode(&config); err != nil {
		fmt.Println("Error decoding config file:", err)
		os.Exit(1)
	}
}

func main() {
	e := echo.New()
	e.HideBanner = true

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format:           `[${time_custom}] ${remote_ip} | ${method} ${path} ${status} | ${latency_human}` + "\n",
		CustomTimeFormat: "2006-01-02 15:04:05",
	}))

	e.GET("/api/:format", func(ctx echo.Context) error {
		timezone := ctx.QueryParam("timezone")
		if timezone == "" {
			timezone = config.DefaultTimezone
		}

		format := ctx.Param("format")
		result, err := progress.Query(format, timezone)
		if err != nil {
			return ctx.JSON(http.StatusBadRequest, map[string]interface{}{"error": err.Error()})
		}
		return ctx.JSONPretty(http.StatusOK, result, "  ")
	})

	if err := e.Start(fmt.Sprintf("%s:%d", config.HostAddress, config.HostPort)); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
