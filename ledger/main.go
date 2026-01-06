package main

import (
	"database/sql"
	"fmt"
	"log"
	"sync"
)

const (
	dbDriver = "mysql"

	dbUser     = "ledger_user"
	dbPassword = "password"

	dbName = "ledger"
	
	topic = "ledger "
)

var wg sync.WaitGroup

func main() {
	dsn := fmt.Sprintf("%s:%s@tcp(localhost:3306)/%s", dbUser, dbPassword, dbName)

	db, err := sql.Open(dbDriver, dsn)
	if err != nil {
		log.Fatalln(err)
	}

	defer func() {
		if err := db.Close(); err != nil {
			log.Println("Error closing ledger db: ", err)
		}
	}()

	if err := db.Ping(); err != nil {
		log.Fatalln(err)
	}
}