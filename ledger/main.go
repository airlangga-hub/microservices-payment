package main

import "sync"

const (
	dbDriver = "mysql"

	dbUser     = "ledger_user"
	dbPassword = "password"

	dbName = "ledger"
	
	topic = "ledger "
)

var wg sync.WaitGroup