package main

type Wallet struct {
	ID         int32  `json:"id"`
	UserID     string `json:"user_id"`
	WalletType string `json:"wallet_type"`
}

type Account struct {
	ID          int32  `json:"id"`
	Cents       int32  `json:"cents"`
	AccountType string `json:"account_type"`
	WalletID    int32  `json:"wallet_id"`
}

type Transaction struct {
	ID                       int32  `json:"id"`
	PID                      string `json:"pid"`
	SrcUserID                string `json:"src_user_id"`
	DstUserID                string `json:"dst_user_id"`
	SrcWalletID              int32  `json:"src_wallet_id"`
	DstWalletID              int32  `json:"dst_wallet_id"`
	SrcAccountID             int32  `json:"src_account_id"`
	DstAccountID             int32  `json:"dst_account_id"`
	SrcAccountType           string `json:"src_account_type"`
	DstAccountType           string `json:"dst_account_type"`
	FinalDstMerchantWalletID int32  `json:"final_dst_merchant_wallet_id"`
	Amount                   int32  `json:"amount"`
}
