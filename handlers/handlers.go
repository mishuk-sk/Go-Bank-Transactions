package handlers

import (
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
type Account struct {
	ID      uuid.UUID `json:"id"`
	UserId  uuid.UUID `json:"user_id"`
	balance float64   `json:"balance"`
}

func Init(router *mux.Router, database *sqlx.DB) {
	db = database

	usersRouter := router.PathPrefix("/users").Subrouter()
	usersRouter.HandleFunc("/", ListUsers).Methods(http.MethodGet)
	usersRouter.HandleFunc("/", CreateUser).Methods(http.MethodPost)
	usersRouter.HandleFunc("/{id}", UpdateUser).Methods(http.MethodPut)
	usersRouter.HandleFunc("/{id}", DeleteUser).Methods(http.MethodDelete)
	usersRouter.HandleFunc("/{id}", GetUser).Methods(http.MethodGet)
	usersRouter.HandleFunc("/{id}/accounts", GetUserAccounts).Methods(http.MethodGet)
	usersRouter.HandleFunc("/{id}/accounts", AddAccount).Methods(http.MethodPost)
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

	user := User{}
	user.ID = uuid.New()
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		raiseErr(err, w, http.StatusBadRequest)
		return
	}
	if _, err := db.Exec("INSERT INTO users(id, first_name, second_name, phone, email) VALUES($1, $2, $3, $4, $5)", user.ID, user.FirstName, user.SecondName, user.Phone, user.Email); err != nil {
		raiseErr(err, w, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	user, err := fetchUser(mux.Vars(r)["id"])
	if err != nil {
		raiseErr(err, w, http.StatusNotFound)
	}
	//TODO deal with id backup
	id := user.ID
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		raiseErr(err, w, http.StatusInternalServerError)
		return
	}
	if _, err := db.Exec("UPDATE users SET first_name = $1, second_name = $2, phone = $3, email = $4 WHERE id=$5", user.FirstName, user.SecondName, user.Phone, user.Email, id); err != nil {
		raiseErr(err, w, http.StatusInternalServerError)
		return
	}
	user.ID = id
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	user, err := fetchUser(mux.Vars(r)["id"])
	if err != nil {
		raiseErr(err, w, http.StatusNotFound)
	}
	if _, err := db.Exec("DELETE FROM users WHERE id=$1", user.ID); err != nil {
		raiseErr(err, w, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	user, err := fetchUser(mux.Vars(r)["id"])
	if err != nil {
		raiseErr(err, w, http.StatusNotFound)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func GetUserAccounts(w http.ResponseWriter, r *http.Request) {
	var accounts []Account
	user, err := fetchUser(mux.Vars(r)["id"])
	if err != nil {
		raiseErr(err, w, http.StatusNotFound)
		return
	}
	if err := db.Select(&accounts, "SELECT * FROM personal_accounts WHERE userId=$1", user.ID); err != nil {
		raiseErr(err, w, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(accounts)
}

func AddAccount(w http.ResponseWriter, r *http.Request) {
	user, err := fetchUser(mux.Vars(r)["id"])
	if err != nil {
		raiseErr(err, w, http.StatusNotFound)
		return
	}
	account := Account{}
	json.NewDecoder(r.Body).Decode(&account)
	account.ID = uuid.New()
	account.UserId = user.ID
	if _, err := db.Exec("INSERT INTO accounts(id, balance, userId) VALUES ($1, $2, $3)", account.ID, account.balance, account.UserId); err != nil {
		raiseErr(err, w, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(account)
}

func fetchUser(id string) (User, error) {
	ID, err := uuid.Parse(id)
	if err != nil {
		return User{}, fmt.Errorf("%s", err.Error())
	}
	user := User{}
	if err := db.Get(&user, "SELECT 1 FROM users WHERE id=$1", ID); err != nil {
		return User{}, fmt.Errorf("%s", err.Error())
	}
	return user, nil
}
