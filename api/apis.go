package api

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	db "github.com/persona-mp3/pwa/database"
	rmq "github.com/rabbitmq/amqp091-go"
	publisher "github.com/persona-mp3/pwa/kafka"
)

var c = publisher.Connection{
	User: "guest",
	Host: "localhost",
	Port: 5672,
	Password: "guest",
}

// TODO: Refactor this to use parameter queries instead
func CreateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	} else if r.Body == nil {
		http.Error(w, "Empty Body", http.StatusBadRequest)
		return
	}

	// now its time to read request body into struct
	// var u *UserReq
	var u *db.UserReq
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	conn, err := db.ConnectDB()
	if err != nil {
		log.Printf("ERROR: createUserHandler : %s", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}

	res, err := conn.CreateUser(u)
	if err != nil {
		log.Printf("ERROR: createUserHandler : %s", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(res)


	// create new connection to kafka
	client, err := publisher.NewConnection(c)
	if err != nil {
		return
	}
	dq := publisher.Queue{
		Name: "break_prod",
		Durable: false,
		AutoDel: false,
		Exclusive: false,
		NoWait: false,
	}
	q, err:= client.DeclareDirectQueue(dq)
	if err != nil {
		return
	}

	
	
	body, err := json.Marshal(res)
	if err != nil {
		log.Println("Could not marshal response??")
		return
	}
	msg := publisher.PublishConfig{
		Exchange: "",
		Key: q.Name, 
		Mandatory: false,
		Immediate: false,
		Msg: &rmq.Publishing{
			ContentType: "text/plain",
			Body: body,
		},
	}
	if err := client.Publish(context.Background(), q, msg); err != nil {
		return
	}

	defer client.Close()

}

