package api

import (
	"encoding/json"
	"log"
	"net/http"

	db "github.com/perona-mp3/pwa/database"
)

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

}
