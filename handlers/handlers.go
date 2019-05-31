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

	validateUserRouter := usersRouter.PathPrefix("/{user_id}/accounts").Subrouter()
	validateUserRouter.Use(checkUserMiddleware)
	validateUserRouter.HandleFunc("/", getUserAccounts).Methods(http.MethodGet)
	validateUserRouter.HandleFunc("/", addAccount).Methods(http.MethodPost)

	// Delete method for accounts was depreciated because of unclear logic under transactions behavior
	// after account disappearing
	//accountsRouter.HandleFunc("/{account_id}/", DeleteAccount).Methods(http.MethodDelete)

	//not validating user account and trnsaction operations
	accountsRouter := usersRouter.PathPrefix("/{user_id}/accounts").Subrouter()
	accountsRouter.HandleFunc("/{account_id}/", getAccount).Methods(http.MethodGet)
	accountsRouter.HandleFunc("/{account_id}/", updateAccount).Methods(http.MethodPut)
	accountsRouter.HandleFunc("/{account_id}/transactions/", getAccountTransactions).Methods(http.MethodGet)

	transactionsRouter := accountsRouter.PathPrefix("/{account_id}/transactions").Subrouter()
	transactionsRouter.Use(checkAccountMiddleware)
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
