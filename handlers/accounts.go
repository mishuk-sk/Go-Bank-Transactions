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
	UserID  uuid.UUID   `json:"user_id" db:"user_id"`
	Name    interface{} `json:"name" db:"name"`
	Balance float64     `json:"balance" db:"balance"`
}

func getAccount(w http.ResponseWriter, r *http.Request) {
	account, err := fetchAccount(mux.Vars(r)["account_id"])
	if err != nil {
		raiseErr(err, w, http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(account)
}

func updateAccount(w http.ResponseWriter, r *http.Request) {
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
	account.Balance = reqAccount.Balance
	account.Name = reqAccount.Name
	if err := updateAcc(account); err != nil {
		raiseErr(err, w, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(account)
}

// Delete method for accounts was depreciated because of unclear logic under transactions behavior
// after account disappearing

/*func DeleteAccount(w http.ResponseWriter, r *http.Request) {
	account, err := fetchAccount(mux.Vars(r)["account_id"])
	if err != nil {
		raiseErr(err, w, http.StatusNotFound)
		return
	}
	tx, err := db.Beginx()
	if err != nil {
		raiseErr(err, w, http.StatusInternalServerError)
		return
	}

	tx.Exec("DELETE FROM transactions WHERE (from_account=$1 AND to_account=NULL) OR (from_account=NULL AND to_account=$1)", account.ID)
	tx.Exec("UPDATE transactions SET from_account=NULL WHERE from_account=$1", account.ID)
	tx.Exec("UPDATE transactions SET to_account=NULL WHERE to_account=$1", account.ID)
	tx.Exec("DELETE FROM personal_accounts WHERE id=$1", account.ID)

	if err := tx.Commit(); err != nil {
		raiseErr(fmt.Errorf("%s, ERROR:%s", "Can't delete account", err.Error()), w, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(account)
}*/
func getUserAccounts(w http.ResponseWriter, r *http.Request) {
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

func addAccount(w http.ResponseWriter, r *http.Request) {
	user, err := fetchUser(mux.Vars(r)["user_id"])
	if err != nil {
		raiseErr(err, w, http.StatusNotFound)
		return
	}
	account := Account{}
	json.NewDecoder(r.Body).Decode(&account)
	account.ID = uuid.New()
	account.UserID = user.ID
	if _, err := db.Exec("INSERT INTO personal_accounts(id, balance, user_id, name) VALUES ($1, $2, $3, $4)", account.ID, account.Balance, account.UserID, account.Name); err != nil {
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

func updateAcc(data Account) error {
	_, err := db.Exec("UPDATE personal_accounts SET balance = $1 name = $2 WHERE id=$2", data.Balance, data.Name, data.ID)
	return err
}
