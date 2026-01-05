package publisher

import (
	"encoding/json"
	"log"
	"time"

	"github.com/IBM/sarama"
)

const (
	emailTopic  = "email"
	ledgerTopic = "ledger"
)

type EmailMessage struct {
	OrderID string `json:"order_id"`
	UserID  string `json:"user_id"`
}

type LedgerMessage struct {
	OrderID   string `json:"order_id"`
	UserID    string `json:"user_id"`
	Amount    int64  `json:"amount"`
	Operation string `json:"operation"`
	Date      string `json:"date"`
}

func SendCaptureMessage(publisher sarama.SyncProducer, pid, srcUserID string, amount int32) {

	emailMessage := EmailMessage{
		OrderID: pid,
		UserID:  srcUserID,
	}

	ledgerMessage := LedgerMessage{
		OrderID: pid,
		UserID:  srcUserID,
		Amount:  int64(amount),
		Date:    time.Now().Format("2006-01-02"),
	}

	sendMessage(publisher, emailMessage, emailTopic)
	sendMessage(publisher, ledgerMessage, ledgerTopic)
}

func sendMessage[T EmailMessage | LedgerMessage](publisher sarama.SyncProducer, message T, topic string) {

	encodedMessage, err := json.Marshal(message)
	if err != nil {
		log.Println("ERROR money movement sendMessage (json.Marshal): ", err)
		return
	}

	kafkaMessage := sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(encodedMessage),
	}

	partition, offset, err := publisher.SendMessage(&kafkaMessage)
	if err != nil {
		log.Println("ERROR money movement sendMessage (publisher.SendMessage): ", err)
		return
	}

	log.Printf("Message sent to partition %d at offset %d\n", partition, offset)
}
