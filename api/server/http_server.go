package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"userService/env"
	"userService/internal/repository"
	"userService/model"
)

type HTTPHandler struct {
	userService repository.UserService
}

func NewHTTPHandler(userService repository.UserService) *HTTPHandler {
	return &HTTPHandler{userService: userService}
}

func (hh *HTTPHandler) StartServer() error {
	r := mux.NewRouter()
	r.HandleFunc("/createUser", hh.createUser).Methods("POST")
	r.HandleFunc("/getUsers", hh.getUsers).Methods("GET")
	r.HandleFunc("/getUser", hh.getUserByID).Methods("GET")

	writeTimeout, errEnv := env.GetTimeDuration("WRITE_TIMEOUT")
	if errEnv != nil {
		return fmt.Errorf("error while getting env vars: %v", errEnv)
	}

	readTimeout, errEnv := env.GetTimeDuration("READ_TIMEOUT")
	if errEnv != nil {
		return fmt.Errorf("error while getting env vars: %v", errEnv)
	}

	server := &http.Server{
		Addr:         os.Getenv("PORT"),
		Handler:      r,
		ReadTimeout:  readTimeout * time.Second,
		WriteTimeout: writeTimeout * time.Second,
	}

	log.Printf("starting server on port 8080")

	err := server.ListenAndServe()
	if err != nil {
		return fmt.Errorf("error starting server on port 8080: %w", err)
	}

	return nil
}

func (hh *HTTPHandler) getUserByID(w http.ResponseWriter, request *http.Request) {
	id := request.URL.Query().Get("id")
	if id == "" {
		log.Print("id is not valid")
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	intID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		log.Printf("error while converting id %v to int: %v", id, err)
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	data, err := hh.userService.GetUserByID(request.Context(), intID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("error getting user by id: %v", err)

		return
	}

	errEncoder := json.NewEncoder(w).Encode(&data)
	if errEncoder != nil {
		log.Printf("error while encoding json: %v", errEncoder)
		w.WriteHeader(http.StatusInternalServerError)

		return
	}
}

func (hh *HTTPHandler) getUsers(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	users, usersErr := hh.userService.GetUsers()
	if usersErr != nil {
		log.Printf("error getting users: %v", usersErr)
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	err := json.NewEncoder(w).Encode(users)
	if err != nil {
		log.Printf("error while encoding json: %v", err)
		w.WriteHeader(http.StatusInternalServerError)

		return
	}
}

func (hh *HTTPHandler) createUser(w http.ResponseWriter, request *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var data repository.Data

	decErr := json.NewDecoder(request.Body).Decode(&data)
	if decErr != nil {
		log.Print(w, decErr)
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	newUser, err := hh.userService.CreateUser(request.Context(), data)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	result := newUser
	createResult := &model.CreateResult{
		Message: "User created successfully",
		Result:  result,
	}

	encoderErr := json.NewEncoder(w).Encode(&createResult)
	if encoderErr != nil {
		log.Print(encoderErr)
		w.WriteHeader(http.StatusInternalServerError)

		return
	}
}
