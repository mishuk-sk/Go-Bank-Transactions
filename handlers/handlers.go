package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/mishuk-sk/Go-Bank-Transactions/subhandler"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
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

	accountsRouter := usersRouter.PathPrefix("/{user_id}/accounts").Subrouter()
	accountsRouter.Use(checkUserMiddleware)
	accountsRouter.HandleFunc("/", GetUserAccounts).Methods(http.MethodGet)
	accountsRouter.HandleFunc("/", AddAccount).Methods(http.MethodPost)
	accountsRouter.HandleFunc("/{account_id}/", GetAccount).Methods(http.MethodGet)
	//accountsRouter.HandleFunc("/{account_id}/", DeleteAccount).Methods(http.MethodDelete)
	accountsRouter.HandleFunc("/{account_id}/", UpdateAccount).Methods(http.MethodPut)

	transactionsRouter := accountsRouter.PathPrefix("/{account_id}/transactions").Subrouter()
	transactionsRouter.Use(checkAccountMiddleware)
	transactionsRouter.HandleFunc("/", GetAccountTransactions).Methods(http.MethodGet)
	// TODO add notify on discarding transaction
	transactionsRouter.HandleFunc("/{transaction_id}/", DiscardTransaction).Methods(http.MethodDelete)
	moneyOperationsRouter := transactionsRouter.PathPrefix("").Subrouter()
	channel := new(subhandler.WorkersChan)
	channel.Init()
	channel.AddListener(notifyUser)
	moneyOperationsRouter.HandleFunc("/", channel.AddWorker(AddTransaction)).Methods(http.MethodPost)
	moneyOperationsRouter.HandleFunc("/enrich/", channel.AddWorker(EnrichAccount)).Methods(http.MethodPost)
	moneyOperationsRouter.HandleFunc("/debit/", channel.AddWorker(DebitAccount)).Methods(http.MethodPost)
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

func checkAccountMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, err := fetchAccount(mux.Vars(r)["account_id"]); err != nil {
			raiseErr(fmt.Errorf("%s; Internal error - %s", "Account not found", err.Error()), w, http.StatusNotFound)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func raiseErr(err error, w http.ResponseWriter, status int) {
	w.WriteHeader(status)
	fmt.Fprintf(w, "%s", err.Error())
	log.Println(err)
}
