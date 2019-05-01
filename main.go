package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

// User type defines user object representation
type User struct {
	ID         string `json:"id"`
	FirstName  string `json:"first_name"`
	SecondName string `json:"second_name"`
	Phone      string `json:"phone"`
	Email      string `json:"email"`
}

func main() {
	router := mux.NewRouter()
	subRouter := router.PathPrefix("/users").Subrouter()

	subRouter.HandleFunc("/", ListUsers).Methods(http.MethodGet)
	subRouter.HandleFunc("/", CreateUser).Methods(http.MethodPost)
}

func ListUsers(w http.ResponseWriter, r *http.Request) {

}

func CreateUser(w http.ResponseWriter, r *http.Request) {

}
