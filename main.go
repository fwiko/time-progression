package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time-progression/progress"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

type config struct {
	HostAddress     string
	HostPort        int
	DefaultTimezone string
}

var cfg config

func init() {
	file, err := os.Open("config.json")
	if err != nil {
		fmt.Println("Error opening config file:", err)
		os.Exit(1)
	}

	decoder := json.NewDecoder(file)
	if err = decoder.Decode(&cfg); err != nil {
		fmt.Println("Error decoding config file:", err)
		os.Exit(1)
	}
}

func main() {
	app := fiber.New(fiber.Config{
		ProxyHeader: "X-Forwarded-For",
	})

	app.Use(logger.New(logger.Config{
		Format:     "[${time}] ${ip} - ${status} ${method} ${path} \n",
		TimeFormat: "2006-01-02 15:04:05",
	}))

	app.Get("/api/:format", func(ctx *fiber.Ctx) error {
		timezone := ctx.Query("timezone")
		if timezone == "" {
			timezone = cfg.DefaultTimezone
		}

		format := ctx.Params("format")
		result, err := progress.Query(format, timezone)
		if err != nil {
			return ctx.Status(400).JSON(fiber.Map{"error": err.Error()})
		}
		return ctx.Status(200).JSON(result)
	})

	if err := app.Listen(fmt.Sprintf("%s:%d", cfg.HostAddress, cfg.HostPort)); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
