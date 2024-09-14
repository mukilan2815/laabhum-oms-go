package nats

import (
	"github.com/nats-io/nats.go"
)

type NatsClient struct {
	conn *nats.Conn
	ec   *nats.EncodedConn
}

func NewNatsClient(url string) *NatsClient {
	nc, _ := nats.Connect(url)
	ec, _ := nats.NewEncodedConn(nc, nats.JSON_ENCODER)

	return &NatsClient{
		conn: nc,
		ec:   ec,
	}
}
