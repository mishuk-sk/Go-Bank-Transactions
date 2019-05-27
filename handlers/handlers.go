package handlers

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/mishuk-sk/Go-Bank-Transactions/workers"
)

type closable interface {
	Close()
}

// FIXME smth went wrong with logic (global variables are widely used, probably should implement "class-like" structure)
var db *sqlx.DB
var toClose []closable

//Init initializes new router with correct routes
func Init(database *sqlx.DB) *mux.Router {

	db = database

	router := mux.NewRouter()
	if val, ok := os.LookupEnv("VERBOSE"); ok && (val == "true") {
		router.Use(verboseMiddleware)
		log.Println("Started in verbose mode")
	}
	usersRouter := router.PathPrefix("/users").Subrouter()
	usersRouter.HandleFunc("/", listUsers).Methods(http.MethodGet)
	usersRouter.HandleFunc("/", createUser).Methods(http.MethodPost)
	usersRouter.HandleFunc("/{user_id}/", updateUser).Methods(http.MethodPut)
	usersRouter.HandleFunc("/{user_id}/", deleteUser).Methods(http.MethodDelete)
	usersRouter.HandleFunc("/{user_id}/", getUser).Methods(http.MethodGet)

	accountsRouter := usersRouter.PathPrefix("/{user_id}/accounts").Subrouter()
	accountsRouter.Use(checkUserMiddleware)
	accountsRouter.HandleFunc("/", getUserAccounts).Methods(http.MethodGet)
	accountsRouter.HandleFunc("/", addAccount).Methods(http.MethodPost)
	accountsRouter.HandleFunc("/{account_id}/", getAccount).Methods(http.MethodGet)
	// Delete method for accounts was depreciated because of unclear logic under transactions behavior
	// after account disappearing

	//accountsRouter.HandleFunc("/{account_id}/", DeleteAccount).Methods(http.MethodDelete)
	accountsRouter.HandleFunc("/{account_id}/", updateAccount).Methods(http.MethodPut)

	transactionsRouter := accountsRouter.PathPrefix("/{account_id}/transactions").Subrouter()
	transactionsRouter.Use(checkAccountMiddleware)
	transactionsRouter.HandleFunc("/", getAccountTransactions).Methods(http.MethodGet)
	// this part implements way to notify listeners on workers pushing to chan
	channel := new(workers.Observer)
	toClose = append(toClose, channel)
	channel.Init()
	channel.AddListener(notifyUser)
	transactionsRouter.HandleFunc("/", channel.CreateHTTPWorker(addTransaction)).Methods(http.MethodPost)
	transactionsRouter.HandleFunc("/enrich/", channel.CreateHTTPWorker(enrichAccount)).Methods(http.MethodPost)
	transactionsRouter.HandleFunc("/debit/", channel.CreateHTTPWorker(debitAccount)).Methods(http.MethodPost)
	transactionsRouter.HandleFunc("/{transaction_id}/", channel.CreateHTTPWorker(discardTransaction)).Methods(http.MethodDelete)
	return router
}

// Close closes everithing initialized by handlers
func Close() {
	for _, v := range toClose {
		v.Close()
	}
}
