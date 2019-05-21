package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type Transaction struct {
	ID          uuid.UUID `json:"id" db:"id"`
	Date        time.Time `json:"date" db:"date"`
	FromAccount uuid.UUID `json:"from_account" db:"from_account"`
	RequestTransaction
}

type RequestTransaction struct {
	ToAccount uuid.UUID `json:"to_account" db:"to_account"`
	Money     float64   `json:"money" db:"money"`
}

func GetAccountTransactions(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(mux.Vars(r)["account_id"])
	if err != nil {
		raiseErr(err, w, http.StatusBadRequest)
		return
	}
	transactions := []Transaction{}
	if err := db.Select(&transactions, "SELECT id, from_account, to_account, date, money::money::numeric::float8 FROM transactions WHERE from_account=$1 OR to_account=$1", id); err != nil {
		raiseErr(err, w, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(transactions)
}

func AddTransaction(w http.ResponseWriter, r *http.Request) {
	from_account, err := fetchAccount(mux.Vars(r)["account_id"])
	if err != nil {
		raiseErr(err, w, http.StatusBadRequest)
		return
	}
	reqTransaction := RequestTransaction{}
	if err := json.NewDecoder(r.Body).Decode(&reqTransaction); err != nil {
		raiseErr(err, w, http.StatusBadRequest)
		return
	}
	if from_account.Balance < reqTransaction.Money {
		raiseErr(fmt.Errorf("%s", "Not enough money on account"), w, http.StatusBadRequest)
		return
	}
	if !checkAccount(reqTransaction.ToAccount) {
		raiseErr(fmt.Errorf("%s", "There is no account exists, that can accept this transaction"), w, http.StatusBadRequest)
		return
	}
	tx, err := db.Beginx()
	if err != nil {
		raiseErr(err, w, http.StatusInternalServerError)
		return
	}
	tx.Exec("UPDATE personal_accounts SET balance = balance - $1 WHERE id=$2", reqTransaction.Money, from_account.ID)
	tx.Exec("UPDATE personal_accounts SET balance = balance + $1 WHERE id=$2", reqTransaction.Money, reqTransaction.ToAccount)
	if err := tx.Commit(); err != nil {
		raiseErr(fmt.Errorf("%s, ERROR:%s", "Can't update accounts", err.Error()), w, http.StatusInternalServerError)
		return
	}
	transaction := Transaction{uuid.New(), time.Now(), from_account.ID, reqTransaction}
	if _, err := db.NamedExec("INSERT INTO transactions(id, date, from_account, to_account, money) VALUES(:id, :date, :from_account, :to_account, :money)", transaction); err != nil {
		raiseErr(fmt.Errorf("%s, ERROR:%s", "Can't create transaction", err.Error()), w, http.StatusInternalServerError)
		tx, err := db.Beginx()
		if err != nil {
			log.Fatalln(err)
			return
		}
		tx.Exec("UPDATE personal_accounts SET balance = balance + $1 WHERE id=$2", reqTransaction.Money, from_account.ID)
		tx.Exec("UPDATE personal_accounts SET balance = balance - $1 WHERE id=$2", reqTransaction.Money, reqTransaction.ToAccount)
		if err := tx.Commit(); err != nil {
			log.Fatalln(err)
			return
		}
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(transaction)
}
func EnrichAccount(w http.ResponseWriter, r *http.Request) {}
func checkAccount(id uuid.UUID) bool {
	var exists bool
	if err := db.QueryRow("SELECT exists (SELECT id FROM personal_accounts WHERE id=$1)", id).Scan(&exists); err != nil && err != sql.ErrNoRows {
		log.Println(err)
		return false
	}
	return exists
}
