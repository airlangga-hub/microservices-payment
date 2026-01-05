package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/IBM/sarama"
)

func AwaitMessages(partitionConsumer sarama.PartitionConsumer, partition int32, done chan struct{}) {

	defer wg.Done()

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

	var em EmailMessage

	if err := json.Unmarshal(msg.Value, &em); err != nil {
		log.Println("INFO error handling message (Unmarshal): ", err)
		return
	}

	if err := SendEmail(em); err != nil {
		log.Println("INFO error handling message (SendEmail): ", err)
		return
	}
}

func SendEmail(em EmailMessage) error {

	return nil
}
