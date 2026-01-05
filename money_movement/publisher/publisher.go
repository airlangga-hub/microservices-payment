package publisher

import (
	"log"
	"sync"
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

func SendCaptureMessage(publisher sarama.AsyncProducer, pid, srcUserID string, amount int32) {

	defer func() {
		if err := publisher.Close(); err != nil {
			log.Println("ERROR money movement SendCaptureMessage (publisher.Close): ", err)
		}
	}()

	emailMessage := EmailMessage{
		OrderID: pid,
		UserID:  srcUserID,
	}

	ledgerMessage := LedgerMessage{
		OrderID: pid,
		UserID:  srcUserID,
		Amount:  int64(amount),
		Date:    time.Now().Format("2020-12-01"),
	}

	
}

func sendMessage[T EmailMessage | LedgerMessage](publisher sarama.AsyncProducer, message T, topic string) {

}
