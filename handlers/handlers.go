package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

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
	moneyOperationsRouter.HandleFunc("/", balanceChange(AddTransaction)).Methods(http.MethodPost)
	moneyOperationsRouter.HandleFunc("/enrich/", balanceChange(EnrichAccount)).Methods(http.MethodPost)
	moneyOperationsRouter.HandleFunc("/debit/", balanceChange(DebitAccount)).Methods(http.MethodPost)
}

// TODO refactor code to reduce request to DB (change calling decorator to work in functions instead, or redesign functions to take db requests as parameters)
// FIXME works incorrect when trying to debit more, than there's on account
func balanceChange(f func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Printf("Error reading body: %v", err)
			http.Error(w, "can't read body", http.StatusBadRequest)
			return
		}
		r.Body = ioutil.NopCloser(bytes.NewBuffer(body))
		f(w, r)
		mainUser, _ := fetchUser(mux.Vars(r)["user_id"])
		mainAccount, _ := fetchAccount(mux.Vars(r)["account_id"])
		reqTransaction := RequestTransaction{}
		json.Unmarshal(body, &reqTransaction)
		log.Print(reqTransaction)
		switch path := r.URL.Path; true {
		case strings.HasSuffix(path, "/enrich/"):
			notifyUser(mainUser, fmt.Sprintf("%s %s. Your account %s(id: %v) was fund with %f", mainUser.FirstName, mainUser.SecondName, mainAccount.Name, mainAccount.ID, reqTransaction.Money))
			return
		case strings.HasSuffix(path, "/debit/"):
			notifyUser(mainUser, fmt.Sprintf("%s %s. Your account %s(id: %v) was debit for %f", mainUser.FirstName, mainUser.SecondName, mainAccount.Name, mainAccount.ID, reqTransaction.Money))
			return
		default:
			notifyUser(mainUser, fmt.Sprintf("%s %s. Your account %s(id: %v) was debit for %f", mainUser.FirstName, mainUser.SecondName, mainAccount.Name, mainAccount.ID, reqTransaction.Money))
			recAccount, _ := fetchAccount(reqTransaction.ToAccount.(string))
			recUser, _ := fetchUser(recAccount.UserID.String())
			notifyUser(recUser, fmt.Sprintf("%s %s. Your account %s(id: %v) was fund with %f", recUser.FirstName, recUser.SecondName, recAccount.Name, recAccount.ID, reqTransaction.Money))
		}

	}
}

// TODO revbuild notifier (email or phone)
func notifyUser(user User, notification string) {
	log.Print(notification)
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
