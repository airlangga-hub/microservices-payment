package main

type LedgerMessage struct {
	OrderID   string `json:"order_id"`
	UserID    string `json:"user_id"`
	Amount    int32  `json:"amount"`
	Operation string `json:"operation"`
	Date      string `json:"date"`
}
