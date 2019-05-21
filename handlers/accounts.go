package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type Account struct {
	ID      uuid.UUID   `json:"id" db:"id"`
	UserId  uuid.UUID   `json:"user_id" db:"user_id"`
	Name    interface{} `json:"name" db:"name"`
	Balance float64     `json:"balance" db:"balance"`
}

func GetAccount(w http.ResponseWriter, r *http.Request) {
	account, err := fetchAccount(mux.Vars(r)["account_id"])
	if err != nil {
		raiseErr(err, w, http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(account)
}

func UpdateAccount(w http.ResponseWriter, r *http.Request) {
	account, err := fetchAccount(mux.Vars(r)["account_id"])
	if err != nil {
		raiseErr(err, w, http.StatusNotFound)
		return
	}
	var reqAccount struct {
		Balance float64 `json:"balance"`
		Name    string  `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&reqAccount); err != nil {
		raiseErr(err, w, http.StatusInternalServerError)
		return
	}
	if _, err := db.Exec("UPDATE personal_accounts SET balance = $1 name = $2 WHERE id=$2", reqAccount.Balance, reqAccount.Name, account.ID); err != nil {
		raiseErr(err, w, http.StatusInternalServerError)
		return
	}
	account.Balance = reqAccount.Balance
	account.Name = reqAccount.Name
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(account)
}

func DeleteAccount(w http.ResponseWriter, r *http.Request) {
	account, err := fetchAccount(mux.Vars(r)["id"])
	if err != nil {
		raiseErr(err, w, http.StatusNotFound)
		return
	}
	if _, err := db.Exec("DELETE FROM personal_accounts WHERE id=$1", account.ID); err != nil {
		raiseErr(err, w, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(account)
}
func GetUserAccounts(w http.ResponseWriter, r *http.Request) {
	var accounts []Account
	user, err := fetchUser(mux.Vars(r)["user_id"])
	if err != nil {
		raiseErr(err, w, http.StatusNotFound)
		return
	}
	if err := db.Select(&accounts, "SELECT user_id, id, name, balance::money::numeric::float8 FROM personal_accounts WHERE user_id=$1", user.ID); err != nil {
		raiseErr(err, w, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(accounts)
}

func AddAccount(w http.ResponseWriter, r *http.Request) {
	user, err := fetchUser(mux.Vars(r)["user_id"])
	if err != nil {
		raiseErr(err, w, http.StatusNotFound)
		return
	}
	account := Account{}
	json.NewDecoder(r.Body).Decode(&account)
	account.ID = uuid.New()
	account.UserId = user.ID
	if _, err := db.Exec("INSERT INTO personal_accounts(id, balance, user_id, name) VALUES ($1, $2, $3, $4)", account.ID, account.Balance, account.UserId, account.Name); err != nil {
		raiseErr(err, w, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(account)
}

func fetchAccount(id string) (Account, error) {
	ID, err := uuid.Parse(id)
	if err != nil {
		return Account{}, fmt.Errorf("%s", err.Error())
	}
	account := Account{}
	if err := db.Get(&account, "SELECT id, name, user_id, balance::money::numeric::float8 FROM personal_accounts WHERE id=$1", ID); err != nil {
		return Account{}, fmt.Errorf("%s", err.Error())
	}
	return account, nil
}
