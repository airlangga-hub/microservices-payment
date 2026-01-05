package main

import (
	"fmt"
	"log"
	"sync"

	"github.com/IBM/sarama"
)

const (
	topic = "email"
)

var wg sync.WaitGroup

func main() {
	done := make(chan struct{})

	consumer, err := sarama.NewConsumer([]string{"kafka:9092"}, sarama.NewConfig())
	if err != nil {
		log.Fatalln("ERROR creating email consumer: ", err)
	}

	defer func() {
		close(done)
		if err := consumer.Close(); err != nil {
			log.Println("ERROR closing email consumer")
		}
	}()

	partitions, err := consumer.Partitions(topic)
	if err != nil {
		log.Fatalln("ERROR getting email partitions: ", err)
	}

	for _, partition := range partitions {
		partitionConsumer, err := consumer.ConsumePartition(topic, partition, sarama.OffsetNewest)
		if err != nil {
			log.Fatalln("ERROR consuming email partition: ", err)
		}

		defer func() {
			if err := partitionConsumer.Close(); err != nil {
				log.Println("ERROR closing email partition consumer: ", err)
			}
		}()

		wg.Add(1)

		go awaitMessages(partitionConsumer, partition, done)
	}

	wg.Wait()
}

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

func handleMessage(msg *sarama.ConsumerMessage)