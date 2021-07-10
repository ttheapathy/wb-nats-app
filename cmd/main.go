package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/nats-io/nats.go"
)

type Messages struct {
	Id   int64  `db:"id"`
	Text string `db:"text"`
}

type ErrorMessage struct {
	Message string
}

func main() {

	app := fiber.New()
	app.Use(logger.New())

	nc := initNats()
	db := initDatabase()

	defer db.Close()
	defer nc.Close()

	for i := 0; i < 3; i++ {
		wrkName := fmt.Sprintf("worker-%d", i+1)
		go worker(nc, db, wrkName)
	}

	var messages []Messages

	for i := 0; i < 1000; i++ {

		messages = append(messages, Messages{Text: "Hello amigos!"})
	}

	app.Post("/insert", func(c *fiber.Ctx) error {

		nc.Publish("insert", messages)

		return c.SendString("👋!")
	})

	app.Get("/get", func(c *fiber.Ctx) error {

		messages := []Messages{}

		if err := db.Select(&messages, "select * from messages order by id desc fetch first 100 rows only"); err != nil {
			return c.JSON(ErrorMessage{Message: err.Error()})
		}

		return c.JSON(messages)
	})

	app.Listen(":3000")
}

func worker(nc *nats.EncodedConn, db *sqlx.DB, name string) {

	log.Printf("Run [%s]", name)

	if _, err := nc.QueueSubscribe("insert", "workers", func(messages []Messages) {
		log.Printf("Started job in: [%s]\n", name)
		_, err := db.NamedExec(`insert into messages (text) VALUES (:text)`, messages)
		if err != nil {
			log.Fatal(err)
		}
	}); err != nil {
		log.Fatal(err)
	}
}

func initDatabase() *sqlx.DB {
	schema := `CREATE TABLE IF NOT EXISTS messages (
		id         serial  not null unique,
		text       varchar(120) NOT NULL
	);`

	dataSourceName := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_SSL"),
	)

	db, err := sqlx.Connect("postgres", dataSourceName)
	if err != nil {
		log.Fatal(err)
	}

	db.MustExec(schema)
	return db
}

func initNats() *nats.EncodedConn {
	nc, err := nats.Connect(os.Getenv("NATS_URL"))
	c, _ := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	if err != nil {
		log.Fatal(err)
	}
	return c
}
