package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

var db *sqlx.DB

func Init(router *mux.Router, database *sqlx.DB) {
	db = database
	usersRouter := router.PathPrefix("/users").Subrouter()
	usersRouter.HandleFunc("/", ListUsers).Methods(http.MethodGet)
	usersRouter.HandleFunc("/", CreateUser).Methods(http.MethodPost)
	usersRouter.HandleFunc("/{id}/", UpdateUser).Methods(http.MethodPut)
	usersRouter.HandleFunc("/{id}/", DeleteUser).Methods(http.MethodDelete)
	usersRouter.HandleFunc("/{id}/", GetUser).Methods(http.MethodGet)
	usersRouter.HandleFunc("/{id}/accounts/", GetUserAccounts).Methods(http.MethodGet)
	usersRouter.HandleFunc("/{id}/accounts/", AddAccount).Methods(http.MethodPost)

	transactionsRouter := usersRouter.PathPrefix("/{userId}").Subrouter()
	transactionsRouter.HandleFunc("/", GetUserTransactions).Methods(http.MethodGet)

	accountsRouter := router.PathPrefix("/accounts/{id}").Subrouter()
	accountsRouter.HandleFunc("/", GetAccount).Methods(http.MethodGet)
	accountsRouter.HandleFunc("/", DeleteAccount).Methods(http.MethodDelete)
	accountsRouter.HandleFunc("/", UpdateAccount).Methods(http.MethodPut)
}

func raiseErr(err error, w http.ResponseWriter, status int) {
	w.WriteHeader(status)
	fmt.Fprintf(w, "%s %s", "Internal server error", err.Error())
	log.Println(err)
}
