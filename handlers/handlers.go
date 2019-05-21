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
	usersRouter.HandleFunc("/{user_id}/", UpdateUser).Methods(http.MethodPut)
	usersRouter.HandleFunc("/{user_id}/", DeleteUser).Methods(http.MethodDelete)
	usersRouter.HandleFunc("/{user_id}/", GetUser).Methods(http.MethodGet)

	accountsRouter := router.PathPrefix("/accounts").Subrouter()
	accountsRouter.Use(checkUserMiddleware)
	accountsRouter.HandleFunc("/{user_id}/", GetUserAccounts).Methods(http.MethodGet)
	usersRouter.HandleFunc("/{user_id}/accounts/", AddAccount).Methods(http.MethodPost)
	accountsRouter.HandleFunc("{account_id}/", GetAccount).Methods(http.MethodGet)
	accountsRouter.HandleFunc("{account_id}/", DeleteAccount).Methods(http.MethodDelete)
	accountsRouter.HandleFunc("{account_id}/", UpdateAccount).Methods(http.MethodPut)

	transactionsRouter := usersRouter.PathPrefix("/{userId}").Subrouter()
	transactionsRouter.HandleFunc("/", GetUserTransactions).Methods(http.MethodGet)

}

func checkUserMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, err := fetchUser(mux.Vars(r)["user_id"]); err != nil {
			raiseErr(fmt.Errorf("%s; Internal error - %s", "User not found", err.Error()), w, http.StatusNotFound)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func raiseErr(err error, w http.ResponseWriter, status int) {
	w.WriteHeader(status)
	fmt.Fprintf(w, "%s %s", "Internal server error", err.Error())
	log.Println(err)
}
