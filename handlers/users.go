package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type User struct {
	ID uuid.UUID `json:"id" db:"id"`
	RequestUser
}
type RequestUser struct {
	FirstName  string      `json:"first_name" db:"first_name"`
	SecondName string      `json:"second_name" db:"second_name"`
	Phone      interface{} `json:"phone" db:"phone"`
	Email      interface{} `json:"email" db:"email"`
}

func ListUsers(w http.ResponseWriter, r *http.Request) {
	var users []User
	err := db.Select(&users, "SELECT * FROM users")
	if err != nil {
		raiseErr(err, w, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(users)
}

func CreateUser(w http.ResponseWriter, r *http.Request) {

	user := User{}
	user.ID = uuid.New()
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		raiseErr(err, w, http.StatusBadRequest)
		return
	}
	if _, err := db.Exec("INSERT INTO users(id, first_name, second_name, phone, email) VALUES($1, $2, $3, $4, $5)", user.ID, user.FirstName, user.SecondName, user.Phone, user.Email); err != nil {
		raiseErr(err, w, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	user, err := fetchUser(mux.Vars(r)["user_id"])
	if err != nil {
		raiseErr(err, w, http.StatusNotFound)
		return
	}
	reqUser := user.RequestUser
	if err := json.NewDecoder(r.Body).Decode(&reqUser); err != nil {
		raiseErr(err, w, http.StatusInternalServerError)
		return
	}
	user.RequestUser = reqUser
	if _, err := db.NamedExec("UPDATE users SET first_name=:first_name, second_name =:second_name, phone =:phone, email =:email WHERE id=:id", user); err != nil {
		raiseErr(err, w, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	user, err := fetchUser(mux.Vars(r)["user_id"])
	if err != nil {
		raiseErr(err, w, http.StatusNotFound)
		return
	}
	if _, err := db.Exec("DELETE FROM users WHERE id=$1", user.ID); err != nil {
		raiseErr(err, w, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	user, err := fetchUser(mux.Vars(r)["user_id"])
	if err != nil {
		raiseErr(err, w, http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
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

func fetchUser(id string) (User, error) {
	ID, err := uuid.Parse(id)
	if err != nil {
		return User{}, fmt.Errorf("%s", err.Error())
	}
	user := User{}
	if err := db.Get(&user, "SELECT * FROM users WHERE id=$1", ID); err != nil {
		return User{}, fmt.Errorf("%s", err.Error())
	}
	return user, nil
}
