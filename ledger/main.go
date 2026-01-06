package main

import (
	"database/sql"
	"fmt"
	"log"
	"sync"

	"github.com/IBM/sarama"
)

const (
	dbDriver = "mysql"

	dbUser     = "ledger_user"
	dbPassword = "password"

	dbName = "ledger"

	topic = "ledger"
)

var (
	wg sync.WaitGroup
	db *sql.DB
)

func main() {
	var err error

	dsn := fmt.Sprintf("%s:%s@tcp(localhost:3306)/%s", dbUser, dbPassword, dbName)

	db, err = sql.Open(dbDriver, dsn)
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

	done := make(chan struct{})

	consumer, err := sarama.NewConsumer([]string{"kafka:9092"}, sarama.NewConfig())
	if err != nil {
		log.Fatalln("ERROR creating ledger consumer: ", err)
	}

	defer func() {
		close(done)
		if err := consumer.Close(); err != nil {
			log.Println("ERROR closing ledger consumer")
		}
	}()

	partitions, err := consumer.Partitions(topic)
	if err != nil {
		log.Fatalln("ERROR getting ledger partitions: ", err)
	}

	for _, partition := range partitions {
		partitionConsumer, err := consumer.ConsumePartition(topic, partition, sarama.OffsetNewest)
		if err != nil {
			log.Fatalln("ERROR opening ledger partition consumer: ", err)
		}

		wg.Add(1)
		go AwaitMessages(partitionConsumer, partition, done)
	}

	wg.Wait()
}
