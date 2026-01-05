package main

import (
	"log"

	"github.com/IBM/sarama"
)

const (
	topic = "email"
)

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

}
