package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func verboseMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		url := r.RequestURI
		method := r.Method
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Printf("Error reading body: %v", err)
			return
		}
		r.Body = ioutil.NopCloser(bytes.NewBuffer(body))
		log.Printf("RequestURI: %s\n Method: %s\n Body: %s\n", url, method, body)
		next.ServeHTTP(w, r)
	})
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
	w.Header().Set("Content-Type", "application/json")
	jsonError := struct {
		Error string
	}{err.Error()}
	json.NewEncoder(w).Encode(jsonError)
	log.Println(err)
}

func notifyUser(not interface{}) {
	notification := not.(Notification)
	user, _ := fetchUser(notification.Account.UserID.String())
	var chargeStr string
	if notification.Debit {
		chargeStr = "charged"
	} else if !notification.Debit {
		chargeStr = "enriched"
	}
	notString := fmt.Sprintf("Dear %s %s, your account %s (id: %v) was %s for %f", user.FirstName, user.SecondName, notification.Account.Name, notification.Account.ID, chargeStr, notification.Transaction.Money)
	if user.Phone != nil {
		log.Printf("SMS: %s\n", notString)
	}
	if user.Email != nil {
		log.Printf("Email: %s\n", notString)
	}
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
