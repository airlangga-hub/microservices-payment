package main

type AuthorizeRequest struct {
	MerchantWalletUserId string `json:"merchant_wallet_user_id"`
	Cents                int64  `json:"cents"`
	Currency             string `json:"currency"`
}

type AuthorizeResponse struct {
	PID string `json:"pid"`
}
