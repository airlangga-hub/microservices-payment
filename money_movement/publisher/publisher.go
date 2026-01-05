package publisher

import (
	"github.com/IBM/sarama"
)

const (
	emailTopic  = "email"
	ledgerTopic = "ledger"
)

func SendCaptureMessage(pid, srcUserID string, amount int32) {

	// create publisher
	publisher, err := sarama.NewSyncProducer()
}
