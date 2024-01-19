package db

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func GetDB() (*sql.DB, error) {
	log.Println("Conectando no banco de dados...")
	db, err := connectSocket()
	if err != nil {
		log.Fatal("erro no arquivo db.go -- database connection error, see .env configuration")
	}
	return db, err
}
