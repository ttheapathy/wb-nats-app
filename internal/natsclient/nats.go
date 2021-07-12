package natsclient

import (
	"log"
	"os"

	"github.com/nats-io/nats.go"
)

func InitNats() *nats.EncodedConn {
	nc, err := nats.Connect(os.Getenv("NATS_URL"))
	c, _ := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	if err != nil {
		log.Fatal(err)
	}
	return c
}
