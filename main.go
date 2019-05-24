package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

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
	db, err := sqlx.Open("postgres", psqlInfo)
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
	// TODO check verbose mode
	router := mux.NewRouter()
	if val, ok := os.LookupEnv("VERBOSE"); ok && (val == "true") {
		router.Use(verboseMiddleware)
		log.Println("Started in verbose mode")
	}
	router.HandleFunc("/", checkLive).Methods(http.MethodGet)
	handlers.Init(router, db)

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
