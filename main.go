package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/persona-mp3/pwa/api"
)

func main() {
	// so this just looks for the index page by default
	fs := http.FileServer(http.Dir("./pwa01/views/"))
	http.HandleFunc("/users/create", api.CreateUser)
	fmt.Println("server running")
	fmt.Println("http//localhost:8700")
	http.Handle("/", fs)
	err := http.ListenAndServe(":8700", nil)
	if err != nil {
		fmt.Println("could not start server")
		log.Fatal(err)
	}
}
