package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/persona-mp3/pwa/api"
)

func main() {
	// so this just looks for the index page by default
	if len(os.Args) <= 1 {
		fmt.Println("Provide port number")
		fmt.Println("go run main.go :8000")
		return
	} else if os.Args[1] == "" {
		fmt.Println("Provide port number")
		fmt.Println("Example: go run main.go :8000")
		return
	}
	// initate connection with RabbitMq before anything else, to avoid opening new connections for every endpoint
	api.RabbitConnect() 

	fs := http.FileServer(http.Dir("./pwa01/views/"))
	http.HandleFunc("/users/create", api.CreateUser)
	fmt.Printf("visit http://localhost%s/users/create\n", os.Args[1])
	http.Handle("/", fs)
	err := http.ListenAndServe(os.Args[1], nil)
	if err != nil {
		fmt.Println("could not start server")
		log.Fatal(err)
	}
}
