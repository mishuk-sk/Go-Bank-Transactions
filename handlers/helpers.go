package handlers

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

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
	fmt.Fprintf(w, "%s", err.Error())
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
