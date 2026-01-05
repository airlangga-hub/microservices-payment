package publisher

const (
	emailTopic  = "email"
	ledgerTopic = "ledger"
)

func SendCaptureMessage(pid, srcUserID string, amount int32)
