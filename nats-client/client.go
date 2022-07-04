package nats_client

import (
	"github.com/nats-io/nats.go"
	"math"
	"time"
)

type NatsClient struct {
	*nats.Conn
}

func NewNatsClient() *NatsClient {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		panic(err)
	}
	return &NatsClient{
		Conn: nc,
	}
}

func (client *NatsClient) Send(topicName string, message string) (data string, err error) {
	msg, err := client.Request(topicName, []byte(message), time.Duration(math.MaxInt64))
	data = string(msg.Data)
	// Simple Publisher
	//return client.Publish(topicName, []byte(message))
	return data, err

}
func (client *NatsClient) SubscribeAsync(topicName string, processFunc func(*nats.Msg)) error {

	// Simple Async Subscriber
	_, err := client.Subscribe(topicName, processFunc)
	return err
}

func (client *NatsClient) Close() {
	// Drain connection (Preferred for responders)
	// Close() not needed if this is called.
	//nc.Drain()
	// Close connection
	client.Close()
}
