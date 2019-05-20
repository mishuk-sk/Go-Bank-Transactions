package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

var db *sqlx.DB

type User struct {
	ID         uuid.UUID   `json:"id" db:"id"`
	FirstName  string      `json:"first_name" db:"first_name"`
	SecondName string      `json:"second_name" db:"second_name"`
	Phone      interface{} `json:"phone" db:"phone"`
	Email      interface{} `json:"email" db:"email"`
}

func Init(router *mux.Router, database *sqlx.DB) {
	db = database

	usersRouter := router.PathPrefix("/users").Subrouter()
	usersRouter.HandleFunc("/", ListUsers).Methods(http.MethodGet)
	usersRouter.HandleFunc("/", CreateUser).Methods(http.MethodPost)
	/*usersRouter.HandleFunc("/{id}", UpdateUser).Methods(http.MethodPut)
	usersRouter.HandleFunc("/{id}", DeleteUser).Methods(http.MethodDelete)
	usersRouter.HandleFunc("/{id}", GetUser).Methods(http.MethodGet)
	usersRouter.HandleFunc("/{id}/accounts", GetUserAccounts).Methods(http.MethodGet)
	*/
}

func raiseErr(err error, w http.ResponseWriter, status int) {
	w.WriteHeader(status)
	fmt.Fprintf(w, "%s %s", "Internal server error", err.Error())
	log.Println(err)
}

func ListUsers(w http.ResponseWriter, r *http.Request) {
	var users []User
	err := db.Select(&users, "SELECT * FROM users")
	if err != nil {
		raiseErr(err, w, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(users)
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	user := User{}
	id := uuid.New()
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		raiseErr(err, w, http.StatusBadRequest)
		return
	}
	rows, err := db.Query("INSERT INTO users(id, first_name, second_name, phone, email) VALUES($1, $2, $3, $4, $5)", id, user.FirstName, user.SecondName, user.Phone, user.Email)
	if err != nil {
		raiseErr(err, w, http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	rows.Next()
	user, _ = scanUser(rows)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

// TODO is it arguable to define scanUser?
func scanUser(row *sql.Rows) (User, error) {
	user := User{}
	if err := row.Scan(&user.ID, &user.FirstName, &user.SecondName, &user.Phone, &user.Email); err != nil {
		return User{}, err
	}
	return user, nil
}
