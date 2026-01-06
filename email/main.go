package main

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/IBM/sarama"
)

const (
	topic = "email"
)

var wg sync.WaitGroup

func main() {
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
		log.Fatalln("ERROR creating email consumer: ", err)
	}

	defer func() {
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
			log.Printf("FATAL: email partition %d failed: %v\n", partition, err)
			once.Do(func() { close(done) })
			wg.Wait()
			return
		}

		wg.Add(1)
		go AwaitMessages(partitionConsumer, partition, done)
	}

	wg.Wait()
}
