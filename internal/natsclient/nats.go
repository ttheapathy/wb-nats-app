package natsclient

import (
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/nats-io/nats.go"
	"github.com/ttheapathy/wb-nats-app/internal/models"
)

func InitNats() *nats.EncodedConn {
	nc, err := nats.Connect(os.Getenv("NATS_URL"))
	c, _ := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	if err != nil {
		log.Fatal(err)
	}
	return c
}

func RunWorker(nc *nats.EncodedConn, db *sqlx.DB, name string) {

	log.Printf("Run [%s]", name)

	if _, err := nc.QueueSubscribe("insert", "workers", func(messages []models.Message) {
		log.Printf("Started job in: [%s]\n", name)
		_, err := db.NamedExec(`insert into messages (text) VALUES (:text)`, messages)
		if err != nil {
			log.Fatal(err)
		}
	}); err != nil {
		log.Fatal(err)
	}
}
