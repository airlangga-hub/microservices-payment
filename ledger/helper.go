package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	"github.com/IBM/sarama"
)

func AwaitMessages(partitionConsumer sarama.PartitionConsumer, partition int32, done chan struct{}) {

	defer func() {
		wg.Done()
		if err := partitionConsumer.Close(); err != nil {
			log.Println("ERROR closing ledger partition consumer: ", err)
		}
	}()

	for {
		select {
		case msg := <-partitionConsumer.Messages():
			fmt.Printf("Partition %d - Received message %s\n", partition, string(msg.Value))
			HandleMessage(msg)
		case <-done:
			fmt.Println("Done signal received, exiting go routine.....")
			return
		}
	}
}

func HandleMessage(msg *sarama.ConsumerMessage) {

	var ledger LedgerMessage

	if err := json.Unmarshal(msg.Value, &ledger); err != nil {
		log.Println("INFO error handling message (Unmarshal): ", err)
		return
	}

	if err := Insert(db, ledger); err != nil {
		log.Println("INFO error handling message (Insert): ", err)
		return
	}
}

func Insert(db *sql.DB, ledger LedgerMessage) error {

	_, err := db.Exec(
		`
		INSERT INTO ledger (
			order_id,
			user_id,
			amount,
			operation,
			date
		) VALUES (
			?,
			?,
			?,
			?,
			?
		)
		`,
		ledger.OrderID,
		ledger.UserID,
		ledger.Amount,
		ledger.Operation,
		ledger.Date,
	)
	if err != nil {
		return err
	}

	return nil
}
