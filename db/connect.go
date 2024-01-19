package db

import (
	"database/sql"
	"fmt"
	"os"
)

func connectSocket() (*sql.DB, error) {
	mustGetenv := func(k string) string {
		v := os.Getenv(k)
		if v == "" {
			return "" // Apenas retorna uma string vazia se a variável de ambiente não estiver configurada
		}
		return v
	}

	var (
		dbUser    = mustGetenv("DB_USER")
		dbPwd     = mustGetenv("DB_PASS")
		dbName    = mustGetenv("DB_NAME")
		dbPort    = mustGetenv("DB_PORT")
		dbTCPHost = mustGetenv("DB_HOST")
	)

	dbURI := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		dbUser, dbPwd, dbTCPHost, dbPort, dbName)

	dbMysql, err := sql.Open("mysql", dbURI)
	if err != nil {
		return nil, err
	}

	return dbMysql, nil
}
