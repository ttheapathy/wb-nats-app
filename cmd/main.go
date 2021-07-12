package main

import (
	"errors"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/nats-io/nats.go"
	"github.com/ttheapathy/wb-nats-app/internal/database"
	"github.com/ttheapathy/wb-nats-app/internal/natsclient"
)

type Message struct {
	Id   int64  `db:"id"`
	Text string `db:"text"`
}

type ErrorMessage struct {
	Message string
	Detail  string
}

func main() {

	app := fiber.New()
	app.Use(logger.New())

	nc := natsclient.InitNats()
	db := database.InitDatabase()

	defer db.Close()
	defer nc.Close()

	for i := 0; i < 3; i++ {
		wrkName := fmt.Sprintf("worker-%d", i+1)
		go worker(nc, db, wrkName)
	}

	var messages []Message

	for i := 0; i < 1000; i++ {

		messages = append(messages, Message{Text: "Hello amigos!"})
	}

	app.Post("/insert", func(c *fiber.Ctx) error {

		nc.Publish("insert", messages)

		return c.SendString("👋!")
	})

	app.Get("/get", func(c *fiber.Ctx) error {

		messages := []Message{}

		if err := db.Select(&messages, "select * from messages order by id desc fetch first 100 rows only"); err != nil {
			return c.JSON(ErrorMessage{Message: errors.New("fetch messages failed").Error(), Detail: err.Error()})
		}

		return c.JSON(messages)
	})

	app.Listen(":3000")
}

func worker(nc *nats.EncodedConn, db *sqlx.DB, name string) {

	log.Printf("Run [%s]", name)

	if _, err := nc.QueueSubscribe("insert", "workers", func(messages []Message) {
		log.Printf("Started job in: [%s]\n", name)
		_, err := db.NamedExec(`insert into messages (text) VALUES (:text)`, messages)
		if err != nil {
			log.Fatal(err)
		}
	}); err != nil {
		log.Fatal(err)
	}
}
