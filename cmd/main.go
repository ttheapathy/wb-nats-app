package main

import (
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	_ "github.com/lib/pq"
	"github.com/ttheapathy/wb-nats-app/internal/database"
	"github.com/ttheapathy/wb-nats-app/internal/models"
	"github.com/ttheapathy/wb-nats-app/internal/natsclient"
)

func main() {

	app := fiber.New()
	app.Use(logger.New())

	nc := natsclient.InitNats()
	db := database.InitDatabase()

	defer db.Close()
	defer nc.Close()

	for i := 0; i < 3; i++ {
		wrkName := fmt.Sprintf("worker-%d", i+1)
		go natsclient.RunWorker(nc, db, wrkName)
	}

	var messages []models.Message

	for i := 0; i < 1000; i++ {

		messages = append(messages, models.Message{Text: "Hello amigos!"})
	}

	app.Post("/insert", func(c *fiber.Ctx) error {

		nc.Publish("insert", messages)

		return c.SendString("👋!")
	})

	app.Get("/get", func(c *fiber.Ctx) error {

		messages := []models.Message{}

		if err := db.Select(&messages, "select * from messages order by id desc fetch first 100 rows only"); err != nil {
			return c.JSON(models.ErrorMessage{Message: errors.New("fetch messages failed").Error(), Detail: err.Error()})
		}

		return c.JSON(messages)
	})

	app.Listen(":3000")
}
