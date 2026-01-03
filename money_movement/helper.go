package main

import "database/sql"

func GetWallet(tx *sql.Tx, userID string) (Wallet, error) {

	var w Wallet

	stmt, err := tx.Prepare(
		`
		SELECT
			id,
			user_id,
			wallet_type
		FROM wallets
		WHERE
			user_id = ?;
		`,
	)
	if err != nil {
		return Wallet{}, err
	}

	if err := stmt.QueryRow(userID).Scan(&w.ID, &w.UserID, &w.WalletType); err != nil {
		return Wallet{}, err
	}

	return w, nil
}

func GetAccount(tx *sql.Tx, walletID int32, accountType string) (Account, error) {

	var a Account

	stmt, err := tx.Prepare(
		`
		SELECT
			id,
			cents,
			account_type,
			wallet_id
		FROM accounts
		WHERE
			wallet_id = ? AND
			account_type = ?;
		`,
	)
	if err != nil {
		return Account{}, err
	}

	if err := stmt.QueryRow(walletID, accountType).Scan(&a.ID, &a.Cents, &a.AccountType, &a.WalletID); err != nil {
		return Account{}, err
	}

	return a, nil
}
