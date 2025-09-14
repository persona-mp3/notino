// Basic configuration for messages to publish to rabbitMq
package rabbit

import (
	"context"
	"fmt"
	"log"

	rmq "github.com/rabbitmq/amqp091-go"
)

const (
	PORT = 5672
	HOST = "localhost"
)


type MqChan struct {
	*rmq.Channel
}

type Queue struct {
	Name      string
	Durable   bool
	AutoDel   bool
	Exclusive bool
	NoWait    bool
}

type PublishConfig struct {
	// Name of the queue initiated from delcaring it
	Exchange *rmq.Queue
	Key       string
	Mandatory bool
	Immediate bool
	Msg       *rmq.Publishing
}

var CURL = fmt.Sprintf("amqp://guest:guest@%s:%d", HOST, PORT)

func errorLogger(s string, err error) {
	if err != nil {
		log.Printf("%s:\n  %s\n", s, err)
	}
}

func NewConnection() (*MqChan, error) {
	conn, err := rmq.Dial(CURL)
	if err != nil {
		errorLogger("panic: could not connect with rabbit", err)
		return nil, err
	}

	// now it's time to create the channel we'll be using
	ch, err := conn.Channel()
	if err != nil {
		errorLogger("error: could not create channel with rabbit", err)
		return nil, err
	}

	// at the moment, since we only have one consumer, we'll use
	// a direct `Exchange` binding
	return &MqChan{ch}, nil
}

func (m *MqChan) DeclareDirectQueue(def Queue) (*rmq.Queue, error) {
	q, err := m.QueueDeclare(
		def.Name, def.Durable, def.AutoDel, def.Exclusive, def.NoWait, nil,
	)
	if err != nil {
		errorLogger(fmt.Sprintf("error declaring queue for: %+v", def), err)
		return nil, err
	}

	// after declaring a queue we can then decide to publish an event
	return &q, nil
}


// Publishes to RabbitMq based on configuration provided in `msg`
//
// Returns errors while trying to publish
func (m *MqChan) Publish(ctx context.Context, q *rmq.Queue, msg PublishConfig) error {
	err := m.PublishWithContext(ctx,
		msg.Exchange.Name,
		msg.Key,
		msg.Mandatory,
		msg.Immediate,
		*msg.Msg,
	)
	if err != nil {
		errorLogger("error occured while trying to publish", err)
		return err
	}
	
	fmt.Println("[o] Message sent")

	return nil
}

