package publisher

import (
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

func SendCaptureMessage(pid, srcUserID string, amount int64) {

	// create publisher
	publisher, err := sarama.NewSyncProducer([]string{"localhost:9092"}, sarama.NewConfig())
	if err != nil {
		log.Println("ERROR money movement SendCaptureMessage (NewSyncProducer): ", err)
	}

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
		Amount:  amount,
		Date:    time.Now().Format("2020-12-01"),
	}

}
