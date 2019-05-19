package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/google/uuid"
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

var db *sql.DB

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
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.DbHost, config.DbPort, config.DbUser, config.DbPassword, config.DbName)

	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	router := mux.NewRouter()
	subRouter := router.PathPrefix("/users").Subrouter()
	router.HandleFunc("/", checkLive).Methods(http.MethodGet)
	subRouter.HandleFunc("/", ListUsers).Methods(http.MethodGet)
	subRouter.HandleFunc("/", CreateUser).Methods(http.MethodPost)
	srv := &http.Server{
		Addr:    "0.0.0.0:8080",
		Handler: router,
	}
	// tests
	log.Printf("%v\n", db)
	log.Printf("%v\n", db.Ping())

	fmt.Println("Server started")
	if err := srv.ListenAndServe(); err != nil {
		log.Println(err)
	}

}

func checkLive(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("LIVE!!!"))
}

func ListUsers(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT * FROM users")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "%s", "Something wrong with db")
		return
	}
	defer rows.Close()
	var users []User
	for rows.Next() {
		usr := User{}
		if err := rows.Scan(&usr.ID, &usr.FirstName, &usr.SecondName, &usr.Phone, &usr.Email); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "%s", "Something wrong with db")
			return
		}
		users = append(users, usr)
	}
	w.WriteHeader(http.StatusInternalServerError)
	marshaledUsers, err := json.Marshal(users)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "%s", "Something wrong with marshalling")
		return
	}
	w.Write(marshaledUsers)
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Internal Server Error")
		return
	}
	id := uuid.New()
	_, err := db.Exec("INSERT INTO users(id, first_name, second_name, phone, email) VALUES ($1, $2, $3, $4, $5)", id, user.FirstName, user.SecondName, user.Phone, user.Email)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "%s", "can't insert into db")
		return
	}
	w.Write([]byte("Inserted"))
}
