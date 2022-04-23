package main

import (
	"authenticator/api"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
)

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/ping", api.Ping).Methods(http.MethodGet)
	router.HandleFunc("/create-user", api.CreateUser).Methods(http.MethodPost)
	router.HandleFunc("/authenticate-user", api.AuthenticateUser).Methods(http.MethodPost)

	log.Println("API is running!")
	_ = http.ListenAndServe(":4000", router)
	defer os.Exit(0)
}
