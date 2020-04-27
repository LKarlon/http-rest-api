package main

import (
	"database/sql"
	"fmt"
	"github.com/LKarlon/http-rest-api.git/api/models"
	_ "github.com/lib/pq"
)

const (
	host     = "tetris.dev.wb.ru"
	port     = 5432
	user     = "postgres"
	password = "popoloka"
	dbname   = "inn_test"
)

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("Successfully connected!")
	db.Exec(
		"INSERT INTO inn_test.inn_shema.ready_data (passport, inn) VALUES ($1, $2)",
		"21 21 91212551",
		"56565446333",
	)
	m := &models.INNReady{}
	db.QueryRow(
		"SELECT passport, inn FROM inn_test.inn_shema.ready_data WHERE passport = $1",
		"21 21 91212551",
	).Scan(&m.Inn, &m.Passport)
	fmt.Println(m)
}

