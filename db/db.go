package db

import (
	"database/sql"
	"go-demoservice/utils"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func MakeDbConnection() (*sql.DB, error) {
	db, err := sql.Open("pgx", utils.ConnectString)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		utils.DbLogger.Error(err)
		return nil, err
	}
	if err := createDbAndTables(db); err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(350)

	return db, nil
}

func createDbAndTables(db *sql.DB) error {
	// if _, err := db.Exec("CREATE DATABASE test_db;"); err != nil {
	// 	return err
	// }

	startQuerysBytes, err := os.ReadFile("db/creation_query.sql")
	if err != nil {
		return err
	}

	startQuerys := string(startQuerysBytes)

	if _, err := db.Exec(startQuerys); err != nil {
		return err
	}

	return nil

}
