package server

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
	"time"
	"userService/internal/user"
)

func StartServer() {
	r := mux.NewRouter()

	r.HandleFunc("/createUser", createUser).Methods("POST")
	r.HandleFunc("/getUsers", getUsers).Methods("GET")
	r.HandleFunc("/getUser", getUserById).Methods("GET")

	server := &http.Server{
		Addr:         ":8080",
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Printf("Starting server on port 8080")

	err := server.ListenAndServe()

	if err != nil {
		log.Fatal(err, "Error starting server on port 8080")
	}
}

func getUserById(w http.ResponseWriter, request *http.Request) {
	id := request.URL.Query().Get("id")

	if id == "" {
		log.Fatal("id is not valid")

		return
	}

	intId, err := strconv.ParseInt(id, 10, 64)

	if err != nil {
		log.Fatal(err)
		return
	}

	data := user.GetUserByID(intId)
	errEncoder := json.NewEncoder(w).Encode(&data)

	if errEncoder != nil {
		log.Fatal(errEncoder)
		return
	}
}

func getUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	err := json.NewEncoder(w).Encode(user.GetUsers())
	if err != nil {
		log.Fatal(err)
		return
	}
}

func createUser(w http.ResponseWriter, request *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	var data user.Data
	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		log.Fatal(w, err)
	}

	result := user.CreateUser(data)
	encoderErr := json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "User created successfully",
		"id":      result.InsertedID,
	})
	if encoderErr != nil {
		log.Fatal(encoderErr)
		return
	}
}
