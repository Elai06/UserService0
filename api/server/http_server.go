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
	"userService/internal/repository"
	"userService/model"
)

type HTTPHandler struct {
	userService repository.UserService
}

func NewHTTPHandler(userService repository.UserService) *HTTPHandler {
	return &HTTPHandler{userService: userService}
}

func (h *HTTPHandler) StartServer(writeTimeout time.Duration, readTimeout time.Duration) error {
	r := mux.NewRouter()
	r.HandleFunc("/createUser", h.createUser).Methods(http.MethodPost)
	r.HandleFunc("/getUsers", h.getUsers).Methods(http.MethodGet)
	r.HandleFunc("/getUser", h.getUserByID).Methods(http.MethodGet)

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

func (h *HTTPHandler) getUserByID(w http.ResponseWriter, request *http.Request) {
	id := request.URL.Query().Get("id")
	if id == "" {
		log.Printf("id is not valid: %v", id)
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	intID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		log.Printf("error while converting id %v to int: %v", id, err)
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	data, err := h.userService.GetUserByID(request.Context(), intID)
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

func (h *HTTPHandler) getUsers(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	users, usersErr := h.userService.GetUsers(req.Context())
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

func (h *HTTPHandler) createUser(w http.ResponseWriter, request *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var data repository.Data
	decErr := json.NewDecoder(request.Body).Decode(&data)
	if decErr != nil {
		log.Printf("error while decoding json: %v", decErr)
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	newUser, err := h.userService.CreateUser(request.Context(), data)
	if err != nil {
		log.Printf("error creating user: %v", err)
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	response := &model.ResponseUserService{
		Message: "User created successfully",
		Result:  newUser,
	}

	encoderErr := json.NewEncoder(w).Encode(&response)
	if encoderErr != nil {
		log.Printf("error while encoding json: %v", encoderErr)
		w.WriteHeader(http.StatusInternalServerError)

		return
	}
}
