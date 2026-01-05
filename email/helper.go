package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/IBM/sarama"
)

func awaitMessages(partitionConsumer sarama.PartitionConsumer, partition int32, done chan struct{}) {

	defer wg.Done()

	for {
		select {
		case msg := <-partitionConsumer.Messages():
			fmt.Printf("Partition %d - Received message %s\n", partition, string(msg.Value))
			handleMessage(msg)
		case <-done:
			fmt.Println("Done signal received, exiting.....")
			return
		}
	}
}

func handleMessage(msg *sarama.ConsumerMessage) {

	var em EmailMessage

	if err := json.Unmarshal(msg.Value, &em); err != nil {
		log.Println("INFO error reading message")
		return
	}
	
	SendEmail(em)
}
