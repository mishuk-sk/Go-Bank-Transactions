package main

import (

	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

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
	var wait time.Duration = 1000000000
	config := configuration{}

	configFile, err := os.Open("config.json")
	if err != nil {
		log.Fatal(err)
	}
	err = json.NewDecoder(configFile).Decode(&config)
	if err != nil {
		log.Fatal(err)
	}
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.DbHost, config.DbPort, config.DbUser, config.DbPassword, config.DbName)
	db, err := sqlx.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
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
	
	router := handlers.Init(db)
	defer handlers.Close()
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


