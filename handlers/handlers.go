package handlers

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/mishuk-sk/Go-Bank-Transactions/workers"
)

type Closable interface {
	Close()
}

// FIXME smth went wrong with logic (global variables are widely used, probably should implement "class-like" structure)
var db *sqlx.DB
var toClose []Closable

func Init(database *sqlx.DB) *mux.Router {

	db = database

	router := mux.NewRouter()
	if val, ok := os.LookupEnv("VERBOSE"); ok && (val == "true") {
		router.Use(verboseMiddleware)
		log.Println("Started in verbose mode")
	}
	usersRouter := router.PathPrefix("/users").Subrouter()
	usersRouter.HandleFunc("/", ListUsers).Methods(http.MethodGet)
	usersRouter.HandleFunc("/", CreateUser).Methods(http.MethodPost)
	usersRouter.HandleFunc("/{user_id}/", UpdateUser).Methods(http.MethodPut)
	usersRouter.HandleFunc("/{user_id}/", DeleteUser).Methods(http.MethodDelete)
	usersRouter.HandleFunc("/{user_id}/", GetUser).Methods(http.MethodGet)

	accountsRouter := usersRouter.PathPrefix("/{user_id}/accounts").Subrouter()
	accountsRouter.Use(checkUserMiddleware)
	accountsRouter.HandleFunc("/", GetUserAccounts).Methods(http.MethodGet)
	accountsRouter.HandleFunc("/", AddAccount).Methods(http.MethodPost)
	accountsRouter.HandleFunc("/{account_id}/", GetAccount).Methods(http.MethodGet)
	// Delete method for accounts was depreciated because of unclear logic under transactions behavior
	// after account disappearing

	//accountsRouter.HandleFunc("/{account_id}/", DeleteAccount).Methods(http.MethodDelete)
	accountsRouter.HandleFunc("/{account_id}/", UpdateAccount).Methods(http.MethodPut)

	transactionsRouter := accountsRouter.PathPrefix("/{account_id}/transactions").Subrouter()
	transactionsRouter.Use(checkAccountMiddleware)
	transactionsRouter.HandleFunc("/", GetAccountTransactions).Methods(http.MethodGet)
	// this part implements way to notify listeners on workers pushing to chan
	channel := new(workers.WorkersChan)
	toClose = append(toClose, channel)
	channel.Init()
	channel.AddListener(notifyUser)
	transactionsRouter.HandleFunc("/", channel.AddHttpWorker(AddTransaction)).Methods(http.MethodPost)
	transactionsRouter.HandleFunc("/enrich/", channel.AddHttpWorker(EnrichAccount)).Methods(http.MethodPost)
	transactionsRouter.HandleFunc("/debit/", channel.AddHttpWorker(DebitAccount)).Methods(http.MethodPost)
	transactionsRouter.HandleFunc("/{transaction_id}/", channel.AddHttpWorker(DiscardTransaction)).Methods(http.MethodDelete)
	return router
}

func Close() {
	for _, v := range toClose {
		v.Close()
	}
}
