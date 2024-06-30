package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/delapaska/cadKeeperAuth/cmd/api"
	"github.com/delapaska/cadKeeperAuth/configs"
	"github.com/delapaska/cadKeeperAuth/db"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		configs.Envs.Host, configs.Envs.DBPort,
		configs.Envs.DBUser, configs.Envs.DBPassword, configs.Envs.DBName)

	db, err := db.NewPostgresSQLStorage(psqlInfo)
	if err != nil {
		log.Fatal(err)
	}

	initStorage(db)
	log.Println("server started on port:", configs.Envs.Port)
	srv := api.NewAPIServer(db)
	srv.Run()

}

func initStorage(db *sql.DB) {
	err := db.Ping()

	if err != nil {
		log.Fatal(err)
	}
	log.Println("DB: Successfully connected")
}
