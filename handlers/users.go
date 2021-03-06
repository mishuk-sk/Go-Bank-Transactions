package handlers

import (
	"encoding/json"
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

func listUsers(w http.ResponseWriter, r *http.Request) {
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

func createUser(w http.ResponseWriter, r *http.Request) {

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

func updateUser(w http.ResponseWriter, r *http.Request) {
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

func deleteUser(w http.ResponseWriter, r *http.Request) {
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

func getUser(w http.ResponseWriter, r *http.Request) {
	user, err := fetchUser(mux.Vars(r)["user_id"])
	if err != nil {
		raiseErr(err, w, http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}
