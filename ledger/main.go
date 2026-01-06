package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

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

	var once sync.Once
	done := make(chan struct{})
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		sig := <-sigChan
		log.Printf("Signal %v received, exiting main....\n", sig)
		once.Do(func() { close(done) })
	}()

	consumer, err := sarama.NewConsumer([]string{"kafka:9092"}, sarama.NewConfig())
	if err != nil {
		log.Fatalln("ERROR creating ledger consumer: ", err)
	}

	defer func() {
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
			log.Printf("FATAL: ledger partition %d failed: %v\n", partition, err)
			once.Do(func() { close(done) })
			wg.Wait()
			return
		}

		wg.Add(1)
		go AwaitMessages(partitionConsumer, partition, done)
	}

	wg.Wait()
}
