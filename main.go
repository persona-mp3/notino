package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/persona-mp3/pwa/api"
	publisher "github.com/persona-mp3/pwa/broker"
)

func main() {
	if len(os.Args) <= 1 {
		fmt.Println("Provide port number ie -> go run main.go :8000 ")
		return
	} else if os.Args[1] == "" {
		fmt.Println("Provide port number ie -> go run main.go :8000 ")
		return
	}

	err := godotenv.Load("./consumer/.env")
	if err != nil {
		fmt.Println("Could not load local varaiables")
		return
	}

	MQ_USER := os.Getenv("MQ_USER")
	MQ_PASSWORD := os.Getenv("MQ_PASSWORD")
	MQ_PORT := os.Getenv("MQ_PORT")
	MQ_HOST := os.Getenv("MQ_HOST")

	i, err := strconv.ParseInt(MQ_PORT, 10, 64)
	if err != nil {
		fmt.Println("Error occured reading port env", err)
		return
	}
	// initate connection with RabbitMq before anything else,
	// to avoid opening new connections for every endpoint
	c := publisher.Connection{
		User:     MQ_USER,
		Password: MQ_PASSWORD,
		Host:     MQ_HOST,
		Port:     i,
	}

	api.RabbitConnect(c)
	fs := http.FileServer(http.Dir("./pwa01/views/"))

	http.HandleFunc("/users/create", api.CreateUser)
	fmt.Printf("visit http://localhost%s/users/create\n", os.Args[1])
	http.Handle("/", fs)

	if err := http.ListenAndServe(os.Args[1], nil); err != nil {
		fmt.Println("could not start server")
		log.Fatal(err)
	}
}
