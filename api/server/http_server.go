package server

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
	"strconv"
	"time"
	"userService/internal/repository"
)

type createResult struct {
	Message string
	Result  *mongo.InsertOneResult
}

type HttpHandler struct {
	userService repository.IUserService
}

func NewHttpHandler(userService repository.IUserService) *HttpHandler {
	return &HttpHandler{userService: userService}
}

func (hh *HttpHandler) StartServer() error {
	r := mux.NewRouter()

	r.HandleFunc("/createUser", hh.createUser).Methods("POST")
	r.HandleFunc("/getUsers", hh.getUsers).Methods("GET")
	r.HandleFunc("/getUser", hh.getUserById).Methods("GET")

	server := &http.Server{
		Addr:         ":8080",
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Printf("Starting server on port 8080")

	err := server.ListenAndServe()

	if err != nil {
		log.Print(err, "Error starting server on port 8080")
		return err
	}

	return nil
}

func (hh *HttpHandler) getUserById(w http.ResponseWriter, request *http.Request) {
	id := request.URL.Query().Get("id")

	if id == "" {
		log.Print("id is not valid")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	intId, err := strconv.ParseInt(id, 10, 64)

	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	data, err := hh.userService.GetUserByID(intId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Print(err)
		return
	}
	errEncoder := json.NewEncoder(w).Encode(&data)

	if errEncoder != nil {
		log.Print(errEncoder)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (hh *HttpHandler) getUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	users, usersErr := hh.userService.GetUsers()
	if usersErr != nil {
		log.Print(usersErr)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err := json.NewEncoder(w).Encode(users)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (hh *HttpHandler) createUser(w http.ResponseWriter, request *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var data repository.Data
	decErr := json.NewDecoder(request.Body).Decode(&data)
	if decErr != nil {
		log.Print(w, decErr)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	newUser, err := hh.userService.CreateUser(data)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	result := newUser
	createResult := &createResult{
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
