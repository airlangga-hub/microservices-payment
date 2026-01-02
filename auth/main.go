package main

import (
	"database/sql"
	"fmt"
	"log"
)

const (
	dbDriver = "mysql"
	
	dbUser = "root"
	dbPassword = "password"
	
	dbName = "users"
)

func main() {
	dsn := fmt.Sprintf("%s:%s@tcp(localhost:3306)/%s", dbUser, dbPassword, dbName)
	
	db, err := sql.Open(dbDriver, dsn)
	if err != nil {
		log.Fatalln(err)
	}
	
	defer func() {
		if err := db.Close(); err != nil {
			log.Println("Error closing db: ", err)
		}
	}()
	
	if err := db.Ping(); err != nil {
		log.Fatalln(err)
	}
}