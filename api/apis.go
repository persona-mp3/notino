// API handlers for aforementioned routes. Each handler, if configured sends,
// publishes messages to RabbitMQ, and database connections. No values are returned
package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	// "github.com/joho/godotenv"

	publisher "github.com/persona-mp3/pwa/broker"
	db "github.com/persona-mp3/pwa/database"
	rmq "github.com/rabbitmq/amqp091-go"
)

var RabbitClient *publisher.Client

func RabbitConnect(c publisher.Connection) {
	conn, err := publisher.NewConnection(c)
	if err != nil {
		log.Printf("PANIC: Could not connect with broker:\n %s\n", err)
		return
	}

	RabbitClient = conn
	log.Println("connection initialised with broker")
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" || r.Body == nil {
		log.Println("client sent a foreign request: ->", r.Method, r.Body)
		http.Error(w, "Bad request\n", http.StatusBadRequest)
		return
	}

	// read request body
	var u *db.UserReq
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, "internal server error", http.StatusTeapot)
		return
	}

	conn, err := db.ConnectDB()
	if err != nil {
		http.Error(w, "internal Server error", http.StatusInternalServerError)
		fmt.Println(err)
		return
	}
	defer conn.Conn.Close()

	// should return a bool value just to certify if the databased refused connection
	// or user data violated database schema
	res, err := conn.CreateUser(u)
	if err != nil {
		log.Printf("ERR: createUserHandler:\n %s\n", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	fmt.Println("sending response to client")
	fmt.Println("response to send -->", res)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(res)

	// now we need to first check if we are still connected with rabit
	if RabbitClient == nil {
		log.Println("PANIC: rabbitClient connection is nil ->", RabbitClient)
		return
	}

	defaultQueue := publisher.Queue{
		Name:      "create_user",
		Durable:   false,
		AutoDel:   false,
		Exclusive: false,
		NoWait:    false,
	}

	declaredQueue, err := RabbitClient.DeclareDirectQueue(defaultQueue)
	if err != nil {
		log.Printf("PANIC: could not declare direct queue:\n %s\n", err)
		return
	}

	payload, err := json.Marshal(res)
	if err != nil {
		log.Println("ERROR: Could not marshal response", err)
		return
	}

	if err := RabbitClient.Ch.PublishWithContext(context.Background(),
		"",
		declaredQueue.Name, false, false,
		rmq.Publishing{
			ContentType: "application/json",
			Body:        []byte(payload),
		},
	); err != nil {
		log.Printf("ERR: Could not publish message to broker:\n %s\n", err)
		return
	}

	log.Println("[o] Message sent to broker")
}
