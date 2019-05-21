package handlers

import (
	"net/http"
)

type Transaction struct {
}

func GetUserTransactions(w http.ResponseWriter, r *http.Request) {
	//user, err := fetchUser(mux.Vars(r)["userId"])
	//if err != nil {
	//	raiseErr(fmt.Errorf("%s\n", "User not found"), w, http.StatusNotFound)
	//	return
	//}
	//db.Select(
}
