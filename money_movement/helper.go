package main

import (
	"database/sql"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

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
	defer stmt.Close()

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
	defer stmt.Close()

	if err := stmt.QueryRow(walletID, accountType).Scan(&a.ID, &a.Cents, &a.AccountType, &a.WalletID); err != nil {
		return Account{}, err
	}

	return a, nil
}

func Transfer(tx *sql.Tx, srcAccount, dstAccount Account, amount int64) error {

	if srcAccount.Cents < int32(amount) {
		return status.Error(codes.InvalidArgument, "not enough money")
	}

	stmt, err := tx.Prepare(
		`
		UPDATE accounts
		SET cents = ?
		WHERE
			id = ?;
		`,
	)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(srcAccount.Cents-int32(amount), srcAccount.ID)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(dstAccount.Cents+int32(amount), dstAccount.ID)
	if err != nil {
		return err
	}

	return nil
}

const (
	insertTransactionQuery = `
		INSERT INTO transactions (
			pid,
			src_user_id,
			dst_user_id,
			src_wallet_id,
			dst_wallet_id,
			src_account_id,
			dst_account_id,
			src_account_type,
			dst_account_type,
			final_dst_merchant_wallet_id,
			amount
		)
		VALUES (
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?
		);
	`

	selectTransactionQuery = `
		SELECT
			id,
			pid,
			src_user_id,
			dst_user_id,
			src_wallet_id,
			dst_wallet_id,
			src_account_id,
			dst_account_id,
			src_account_type,
			dst_account_type,
			final_dst_merchant_wallet_id,
			amount
		FROM transactions
		WHERE
			pid = ?;
	`
)

func CreateTransaction(tx *sql.Tx, srcAccount, dstAccount Account, srcUserID, dstUserID string, merchantWalletID int32, amount int64) (string, error) {

	pid := uuid.NewString()

	stmt, err := tx.Prepare(insertTransactionQuery)
	if err != nil {
		return "", err
	}
	defer stmt.Close()

	_, err = stmt.Exec(
		pid,
		srcUserID,
		dstUserID,
		srcAccount.WalletID,
		dstAccount.WalletID,
		srcAccount.ID,
		dstAccount.ID,
		srcAccount.AccountType,
		dstAccount.AccountType,
		merchantWalletID,
		int32(amount),
	)
	if err != nil {
		return "", err
	}

	return pid, nil
}

func GetTransaction(tx *sql.Tx, pid string) (Transaction, error) {

	var t Transaction

	stmt, err := tx.Prepare(selectTransactionQuery)
	if err != nil {
		return Transaction{}, err
	}
	defer stmt.Close()

	if err := stmt.QueryRow(pid).Scan(
		&t.ID,
		&t.PID,
		&t.SrcUserID,
		&t.DstUserID,
		&t.SrcWalletID,
		&t.DstWalletID,
		&t.SrcAccountID,
		&t.DstAccountID,
		&t.SrcAccountType,
		&t.DstAccountType,
		&t.FinalDstMerchantWalletID,
		&t.Amount,
	); err != nil {
		return Transaction{}, err
	}

	return t, nil
}

func GetWalletByID(tx *sql.Tx, walletID int32) (Wallet, error) {
    var w Wallet
    query := `SELECT id, user_id FROM wallets WHERE id = ?`
    err := tx.QueryRow(query, walletID).Scan(&w.ID, &w.UserID)
    return w, err
}