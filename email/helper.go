package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/smtp"

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

	senderEmail := "your_email@gmail.com"
	password := "your_password"

	message := fmt.Appendf(nil, "Subject: Payment Processed!\nProcess ID: %s", em.OrderID)

	smtpServer := "smtp.gmail.com"
	smtpPort := 587

	creds := smtp.PlainAuth("", senderEmail, password, smtpServer)

	smtpAddress := fmt.Sprintf("%s:%d", smtpServer, smtpPort)

	if err := smtp.SendMail(smtpAddress, creds, senderEmail, []string{em.UserID}, message); err != nil {
		return err
	}

	return nil
}
