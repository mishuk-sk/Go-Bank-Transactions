package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type configuration struct {
	DbName     string `json:"db_name"`
	DbPort     int    `json:"db_port"`
	DbHost     string `json:"db_host"`
	DbUser     string `json:"db_user"`
	DbPassword string `json:"db_password"`
}

// User type defines user object representation
type User struct {
	ID         string `json:"id"`
	FirstName  string `json:"first_name"`
	SecondName string `json:"second_name"`
	Phone      string `json:"phone"`
	Email      string `json:"email"`
}

func main() {
	configFile, err := os.Open("config.json")
	if err != nil {
		panic(err)
	}
	config := configuration{}
	err = json.NewDecoder(configFile).Decode(&config)
	if err != nil {
		panic(err)
	}
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s",
		config.DbHost, config.DbPort, config.DbUser, config.DbPassword, config.DbName)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	router := mux.NewRouter()
	subRouter := router.PathPrefix("/users").Subrouter()

	subRouter.HandleFunc("/", ListUsers).Methods(http.MethodGet)
	subRouter.HandleFunc("/", CreateUser).Methods(http.MethodPost)
}

func ListUsers(w http.ResponseWriter, r *http.Request) {

}

func CreateUser(w http.ResponseWriter, r *http.Request) {

}
