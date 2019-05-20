package main

// TODO create different packages
import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/mishuk-sk/Go-Bank-Transactions/handlers"
)

type configuration struct {
	DbName     string `json:"db_name"`
	DbPort     int    `json:"db_port"`
	DbHost     string `json:"db_host"`
	DbUser     string `json:"db_user"`
	DbPassword string `json:"db_password"`
}
type User struct {
	ID         uuid.UUID   `json:"id"`
	FirstName  string      `json:"first_name"`
	SecondName string      `json:"second_name"`
	Phone      interface{} `json:"phone"`
	Email      interface{} `json:"email"`
}

type Account struct {
	ID      uuid.UUID `json:"id"`
	UserId  uuid.UUID `json:"user_id"`
	balance float64   `json:"balance"`
}

var db *sqlx.DB

func main() {
	var wait time.Duration
	config := configuration{}

	configFile, err := os.Open("config.json")
	if err != nil {
		panic(err)
	}
	err = json.NewDecoder(configFile).Decode(&config)
	if err != nil {
		panic(err)
	}
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.DbHost, config.DbPort, config.DbUser, config.DbPassword, config.DbName)
	db, err = sqlx.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	// test db connection
	if err := db.Ping(); err != nil {
		log.Printf("%v\n", err)
	} else {
		log.Println("Db successfully connected")
	}

	// initializing routes
	// TODO add vendoring
	router := mux.NewRouter()
	router.HandleFunc("/", checkLive).Methods(http.MethodGet)
	handlers.Init(router, db)
	/*usersRouter := router.PathPrefix("/users").Subrouter()
	usersRouter.HandleFunc("/", ListUsers).Methods(http.MethodGet)
	usersRouter.HandleFunc("/", CreateUser).Methods(http.MethodPost)
	usersRouter.HandleFunc("/{id}", UpdateUser).Methods(http.MethodPut)
	usersRouter.HandleFunc("/{id}", DeleteUser).Methods(http.MethodDelete)
	usersRouter.HandleFunc("/{id}", GetUser).Methods(http.MethodGet)
	usersRouter.HandleFunc("/{id}/accounts", GetUserAccounts).Methods(http.MethodGet)
	*/
	//usersRouter.HandleFunc("/{id}/accounts", AddAccount).Methods(http.MethodPost)
	// TODO add transactions and personal accounts
	// http server
	srv := &http.Server{
		Addr:    "0.0.0.0:8080",
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	srv.Shutdown(ctx)
	log.Println("shutting down")
	os.Exit(0)

}

func checkLive(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("LIVE!!!"))
}

/*func GetUserAccounts(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err.Error())
		log.Println(err)
		return
	}
	exists := false
	db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE id=$1", id).Scan(&exists)
	if !exists {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "%s\n", "User: $1 not found", id)
		log.Printf("%s\n", "User: $1 not found", id)
		return
	}

	rows, err := db.Query("SELECT * FROM personal_accounts WHERE userId=$1", id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err)
		log.Println(err)
		return
	}
	defer rows.Close()
	var accounts []Account
	for rows.Next() {
		account := Account{}
		if err := rows.Scan(&account.ID, &account.balance, &account.UserId); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, err)
			log.Println(err)
			return
		}
		accounts = append(accounts, account)
	}
	marshaledAccounts, err := json.Marshal(accounts)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err)
		log.Println(err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(marshaledAccounts)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err.Error())
		log.Println(err)
		return
	}
	usr, err := getUser(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, err.Error())
		log.Println(err)
		return
	}
	if err := json.NewDecoder(r.Body).Decode(&usr); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err.Error())
		log.Println(err)
		return
	}
	if _, err := db.Exec("UPDATE users SET first_name = $1, second_name = $2, phone = $3, email = $4 WHERE id=$5", usr.FirstName, usr.SecondName, usr.Phone, usr.Email, usr.ID); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err.Error())
		log.Println(err)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Successfully updated")

}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err.Error())
		log.Println(err)
		return
	}
	if _, err := db.Exec("DELETE FROM users WHERE id=$1", id); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err.Error())
		log.Println(err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err.Error())
		log.Println(err)
		return
	}
	user, err := getUser(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, err.Error())
		log.Println(err)
		return
	}
	marshaledUser, err := json.Marshal(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err.Error())
		log.Println(err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(marshaledUser)
}

func getUser(id uuid.UUID) (User, error) {
	user := User{}
	row := db.QueryRow("SELECT * FROM users WHERE id=$1", id)
	if err := row.Scan(&user.ID, &user.FirstName, &user.SecondName, &user.Phone, &user.Email); err != nil {
		return user, fmt.Errorf("%s\n", "User doesn't exist")
	}
	return user, nil
}

func ListUsers(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT * FROM users")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "%s", "Can't connect to db")
		log.Println(err)
		return
	}
	defer rows.Close()
	var users []User
	for rows.Next() {
		usr := User{}
		if err := rows.Scan(&usr.ID, &usr.FirstName, &usr.SecondName, &usr.Phone, &usr.Email); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "%s", "problem scanning row")
			log.Println(err)
			return
		}
		users = append(users, usr)
	}

	marshaledUsers, err := json.Marshal(users)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "%s", "Something wrong with marshalling")
		log.Println(err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(marshaledUsers)
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	user := User{}
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Wrong payload format")
		log.Println(err)
		return
	}
	user.ID = uuid.New()

	if _, err := db.Exec("INSERT INTO users(id, first_name, second_name, phone, email) VALUES ($1, $2, $3, $4, $5)", user.ID, user.FirstName, user.SecondName, user.Phone, user.Email); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "%s", "Can't insert into db")
		log.Println(err)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s %v\n", "Inserted", user)
}
*/
