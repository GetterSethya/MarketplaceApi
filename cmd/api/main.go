package main

import (
	"log"
	"time"

	"github.com/GetterSethya/golangApiMarketplace/config"
	"github.com/GetterSethya/golangApiMarketplace/internal/datastore"
	"github.com/GetterSethya/golangApiMarketplace/internal/server"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	cfg := config.LoadConfig()
	sqlStorage := datastore.NewPostgresStorage()

	db, err := sqlStorage.Init()
	if err != nil {
		log.Fatal(err)
	}

	// set db conn limit
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	store := datastore.NewStore(db)
	api := server.NewServer(cfg.App.Port, store)

	api.Run()
}
