package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type Transaction struct {
	ID          uuid.UUID   `json:"id" db:"id"`
	Date        time.Time   `json:"date" db:"date"`
	FromAccount interface{} `json:"from_account" db:"from_account"`
	RequestTransaction
}

type RequestTransaction struct {
	ToAccount interface{} `json:"to_account" db:"to_account"`
	Money     float64     `json:"money" db:"money"`
}

type Notification struct {
	Account     Account
	Transaction Transaction
	Debit       bool
}

func getAccountTransactions(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(mux.Vars(r)["account_id"])
	if err != nil {
		raiseErr(err, w, http.StatusBadRequest)
		return
	}
	var transactions []struct {
		ID          uuid.UUID `json:"id" db:"id"`
		Date        time.Time `json:"date" db:"date"`
		FromAccount uuid.UUID `json:"from_account" db:"from_account"`
		ToAccount   uuid.UUID `json:"to_account" db:"to_account"`
		Money       float64   `json:"money" db:"money"`
	}
	if err := db.Select(&transactions, "SELECT id, from_account, to_account, date, money::money::numeric::float8 FROM transactions WHERE from_account=$1 OR to_account=$1", id); err != nil {
		raiseErr(err, w, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(transactions)
}

func addTransaction(w http.ResponseWriter, r *http.Request, ch chan<- interface{}) {
	fromAccount, _ := fetchAccount(mux.Vars(r)["account_id"])
	reqTransaction := RequestTransaction{}
	if err := json.NewDecoder(r.Body).Decode(&reqTransaction); err != nil {
		raiseErr(err, w, http.StatusBadRequest)
		return
	}
	if fromAccount.Balance < reqTransaction.Money {
		raiseErr(fmt.Errorf("%s", "Not enough money on account"), w, http.StatusBadRequest)
		return
	}
	toAccount, err := fetchAccount(reqTransaction.ToAccount.(string))
	if err != nil {
		raiseErr(fmt.Errorf("%s", "There is no account exists, that can accept this transaction"), w, http.StatusBadRequest)
		return
	}
	transaction := Transaction{uuid.New(), time.Now(), fromAccount.ID, reqTransaction}
	tx, err := db.Beginx()
	if err != nil {
		raiseErr(err, w, http.StatusInternalServerError)
		return
	}
	tx.Exec("UPDATE personal_accounts SET balance = balance - $1 WHERE id=$2", reqTransaction.Money, fromAccount.ID)
	tx.Exec("UPDATE personal_accounts SET balance = balance + $1 WHERE id=$2", reqTransaction.Money, reqTransaction.ToAccount)
	tx.NamedExec("INSERT INTO transactions(id, date, from_account, to_account, money) VALUES(:id, :date, :from_account, :to_account, :money)", transaction)
	if err := tx.Commit(); err != nil {
		raiseErr(fmt.Errorf("%s, ERROR:%s", "Can't create transaction", err.Error()), w, http.StatusInternalServerError)
		return
	}
	// Making balance relevant for both account after transaction
	fromAccount.Balance -= transaction.Money
	toAccount.Balance += transaction.Money
	ch <- Notification{fromAccount, transaction, true}
	ch <- Notification{toAccount, transaction, false}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(transaction)
}

func enrichAccount(w http.ResponseWriter, r *http.Request, ch chan<- interface{}) {
	req := RequestTransaction{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		raiseErr(err, w, http.StatusBadRequest)
		return
	}
	account, err := fetchAccount(mux.Vars(r)["account_id"])
	if err != nil {
		raiseErr(err, w, http.StatusBadRequest)
		return
	}
	req.ToAccount = account.ID
	transaction := Transaction{uuid.New(), time.Now(), nil, req}
	tx, err := db.Beginx()
	if err != nil {
		raiseErr(err, w, http.StatusInternalServerError)
		return
	}
	tx.Exec("UPDATE personal_accounts SET balance = balance + $1 WHERE id=$2", req.Money, req.ToAccount)
	tx.NamedExec("INSERT INTO transactions(id, date, from_account, to_account, money) VALUES(:id, :date, :from_account, :to_account, :money)", transaction)
	if err := tx.Commit(); err != nil {
		raiseErr(fmt.Errorf("%s, ERROR:%s", "Can't create transaction", err.Error()), w, http.StatusInternalServerError)
		return
	}
	// Making balance relevant after transaction
	account.Balance += transaction.Money
	ch <- Notification{account, transaction, false}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(transaction)
}

func debitAccount(w http.ResponseWriter, r *http.Request, ch chan<- interface{}) {
	req := RequestTransaction{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		raiseErr(err, w, http.StatusBadRequest)
		return
	}
	fromAccount, err := fetchAccount(mux.Vars(r)["account_id"])
	if err != nil {
		raiseErr(err, w, http.StatusBadRequest)
		return
	}
	if fromAccount.Balance < req.Money {
		raiseErr(fmt.Errorf("%s", "Not enough money on account"), w, http.StatusBadRequest)
		return
	}
	req.ToAccount = nil
	transaction := Transaction{uuid.New(), time.Now(), fromAccount.ID, req}
	tx, err := db.Beginx()
	if err != nil {
		raiseErr(err, w, http.StatusInternalServerError)
		return
	}
	tx.Exec("UPDATE personal_accounts SET balance = balance - $1 WHERE id=$2", req.Money, fromAccount.ID)
	tx.NamedExec("INSERT INTO transactions(id, date, from_account, to_account, money) VALUES(:id, :date, :from_account, :to_account, :money)", transaction)
	if err := tx.Commit(); err != nil {
		raiseErr(fmt.Errorf("%s, ERROR:%s", "Can't create transaction", err.Error()), w, http.StatusInternalServerError)
		return
	}
	// Making balance relevant after transaction
	fromAccount.Balance -= transaction.Money
	ch <- Notification{fromAccount, transaction, true}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(transaction)
}

func discardTransaction(w http.ResponseWriter, r *http.Request, ch chan<- interface{}) {
	// Defining new struct for transaction because of incorrect parsing with regular Transaction type
	// (Exactly UUIDs are parsed as bytes and need to be converted after parsing)
	var transaction struct {
		ID          uuid.UUID `json:"id" db:"id"`
		Date        time.Time `json:"date" db:"date"`
		FromAccount uuid.UUID `json:"from_account" db:"from_account"`
		ToAccount   uuid.UUID `json:"to_account" db:"to_account"`
		Money       float64   `json:"money" db:"money"`
	}
	id, err := uuid.Parse(mux.Vars(r)["transaction_id"])
	if err != nil {
		raiseErr(fmt.Errorf("%s: %s", "Wrong transaction id", err.Error()), w, http.StatusBadRequest)
		return
	}
	if err := db.Get(&transaction, "SELECT id, date, from_account, to_account, money::money::numeric::float8 FROM transactions WHERE id=$1", id); err != nil {
		raiseErr(fmt.Errorf("%s: %s", "No transaction exists with this id", err.Error()), w, http.StatusNotFound)
		return
	}
	tx, err := db.Beginx()
	if err != nil {
		raiseErr(err, w, http.StatusInternalServerError)
		return
	}
	if transaction.FromAccount != uuid.Nil {
		tx.Exec("UPDATE personal_accounts SET balance = balance + $1 WHERE id=$2", transaction.Money, transaction.FromAccount)
	}
	if transaction.ToAccount != uuid.Nil {
		tx.Exec("UPDATE personal_accounts SET balance = balance - $1 WHERE id=$2", transaction.Money, transaction.ToAccount)
	}
	tx.Exec("DELETE FROM transactions WHERE id=$1", transaction.ID)
	if err := tx.Commit(); err != nil {
		raiseErr(fmt.Errorf("%s, ERROR:%s", "Can't delete transaction", err.Error()), w, http.StatusInternalServerError)
		return
	}
	tr := Transaction{transaction.ID, transaction.Date, transaction.FromAccount, RequestTransaction{transaction.ToAccount, transaction.Money}}
	if transaction.FromAccount != uuid.Nil {
		fromAccount, err := fetchAccount(transaction.FromAccount.String())
		if err != nil {
			raiseErr(fmt.Errorf("Can't fetch account %v to notify user", transaction.FromAccount), w, http.StatusInternalServerError)
		}
		ch <- Notification{fromAccount, tr, false}
	}
	if transaction.ToAccount != uuid.Nil {
		toAccount, err := fetchAccount(transaction.ToAccount.String())
		if err != nil {
			raiseErr(fmt.Errorf("Can't fetch account %v to notify user", transaction.ToAccount), w, http.StatusInternalServerError)
		}
		ch <- Notification{toAccount, tr, true}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(transaction)
}
