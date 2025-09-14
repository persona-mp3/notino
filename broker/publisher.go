// Basic configuration for messages to publish to rabbitMq
package publisher

import (
	"context"
	"fmt"
	"log"

	rmq "github.com/rabbitmq/amqp091-go"
)

type Connection struct {
	User     string
	Host     string
	Password string
	Port     int64
}
type Client struct {
	conn *rmq.Connection
	ch   *rmq.Channel
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
	Exchange  string
	Key       string
	Mandatory bool
	Immediate bool
	Msg       *rmq.Publishing
}

func errorLogger(s string, err error) {
	if err != nil {
		log.Printf("%s:\n  %s\n", s, err)
	}
}

func NewConnection(c Connection) (*Client, error) {
	// var CURL = fmt.Sprintf("amqp://%s:%s@%s:%d", c.User, c.Password, c.Host, c.Port)
	var CURL = fmt.Sprintf("amqp://guest:guest@%s:%d", c.Host, c.Port)
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
	return &Client{ch: ch, conn: conn}, nil
}

func (c *Client) Close() error {
	if err := c.conn.Close(); err != nil {
		errorLogger("error in closing connection", err)
		return err
	}

	if err := c.ch.Close(); err != nil {
		errorLogger("error in closing channel", err)
		return err
	}
	return nil
}

func (c *Client) DeclareDirectQueue(def Queue) (*rmq.Queue, error) {
	q, err := c.ch.QueueDeclare(
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
func (c *Client) Publish(ctx context.Context, q *rmq.Queue, msg PublishConfig) error {
	err := c.ch.PublishWithContext(ctx,
		msg.Exchange,
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

// I'm not sure if we should make a seperae consumer file, but since it doesn't need much 
// it'd be best to just leave it here in the mean-time just incase


type ConsumeConfig struct {
	// This should be the same as the name for the publish queue since 
	// it's using the default exchange method. since it binds it to the route
	Name string
}
func (c *Client) Consume(queue *rmq.Queue) error {
	msgs, err := c.ch.Consume(queue.Name, "", false, false, false, false, nil)
	if err != nil {
		errorLogger("Could not consume", err)
		return err
	}

	var forever = make(chan bool)
	go func(){
		for d := range msgs {
			fmt.Printf("[*]New Notification\n\n")
			log.Printf("\n%s\n", d.Body)
		}
	}()

	// Block main goroutine to listen for incoming request, by polling??
	fmt.Println("**Waiting for new messages**")
	<-forever
	return nil
}
