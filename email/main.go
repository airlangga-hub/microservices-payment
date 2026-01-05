package main

import (
	"log"

	"github.com/IBM/sarama"
)

const (
	topic = "email"
)

func main() {

	consumer, err := sarama.NewConsumer()
	if err != nil {
		log.Fatalln("ERROR creating email consumer: ", err)
	}

}
