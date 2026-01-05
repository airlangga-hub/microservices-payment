package main

import (
	"context"
	"database/sql"
	"log"

	"github.com/IBM/sarama"
	"github.com/airlangga-hub/microservices-payment/money_movement/pb"
	"github.com/airlangga-hub/microservices-payment/money_movement/publisher"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Server struct {
	pb.UnimplementedMoneyMovementServiceServer
	db        *sql.DB
	publisher sarama.AsyncProducer
}

func NewServer(db *sql.DB, publisher sarama.AsyncProducer) *Server {
	return &Server{db: db, publisher: publisher}
}

func (s *Server) Authorize(ctx context.Context, r *pb.AuthorizeRequest) (*pb.AuthorizeResponse, error) {

	if r.Currency != "USD" {
		return nil, status.Error(codes.InvalidArgument, "only accepts USD")
	}

	// begin tx
	tx, err := s.db.Begin()
	if err != nil {
		log.Println("ERROR money movement Authorize (db.Begin): ", err)
		return nil, status.Error(codes.Internal, "failed to authorize transaction")
	}

	defer tx.Rollback()

	merchantWallet, err := GetWallet(tx, r.MerchantWalletUserId)
	if err != nil {
		log.Println("ERROR money movement Authorize (GetWallet): ", err)
		return nil, status.Error(codes.NotFound, "invalid merchant wallet user id")
	}

	customerWallet, err := GetWallet(tx, r.CustomerWalletUserId)
	if err != nil {
		log.Println("ERROR money movement Authorize (GetWallet): ", err)
		return nil, status.Error(codes.NotFound, "invalid customer wallet user id")
	}

	srcAccount, err := GetAccount(tx, customerWallet.ID, "DEFAULT")
	if err != nil {
		log.Println("ERROR money movement Authorize (GetWallet): ", err)
		return nil, status.Error(codes.NotFound, "source account not found")
	}

	dstAccount, err := GetAccount(tx, customerWallet.ID, "PAYMENT")
	if err != nil {
		log.Println("ERROR money movement Authorize (GetWallet): ", err)
		return nil, status.Error(codes.NotFound, "destination account not found")
	}

	err = Transfer(tx, srcAccount, dstAccount, r.Cents)
	if err != nil {
		log.Println("ERROR money movement Authorize (Transfer): ", err)
		return nil, status.Error(codes.Internal, "failed transfer")
	}

	pid := uuid.NewString()

	err = CreateTransaction(tx, pid, srcAccount, dstAccount, customerWallet.UserID, customerWallet.UserID, merchantWallet.ID, r.Cents)
	if err != nil {
		log.Println("ERROR money movement Authorize (CreateTransaction): ", err)
		return nil, status.Error(codes.Internal, "failed creating transaction")
	}

	// end tx
	err = tx.Commit()
	if err != nil {
		log.Println("ERROR money movement Authorize (tx.Commit): ", err)
		return nil, status.Error(codes.Internal, "failed commiting transaction")
	}

	return &pb.AuthorizeResponse{Pid: pid}, nil
}

func (s *Server) Capture(ctx context.Context, r *pb.CaptureRequest) (*emptypb.Empty, error) {

	// begin tx
	tx, err := s.db.Begin()
	if err != nil {
		log.Println("ERROR money movement Capture (db.Begin): ", err)
		return nil, status.Error(codes.Internal, "failed to capture transaction")
	}

	defer tx.Rollback()

	authorizedTransaction, err := GetTransaction(tx, r.Pid)
	if err != nil {
		log.Println("ERROR money movement Capture (GetTransaction): ", err)
		return nil, status.Error(codes.Internal, "transaction not found")
	}

	srcAccount, err := GetAccount(tx, authorizedTransaction.DstWalletID, "PAYMENT")
	if err != nil {
		log.Println("ERROR money movement Capture (GetAccount): ", err)
		return nil, status.Error(codes.Internal, "source account not found")
	}

	dstMerchantAccount, err := GetAccount(tx, authorizedTransaction.FinalDstMerchantWalletID, "INCOMING")
	if err != nil {
		log.Println("ERROR money movement Capture (GetAccount): ", err)
		return nil, status.Error(codes.Internal, "destination account not found")
	}

	err = Transfer(tx, srcAccount, dstMerchantAccount, int64(authorizedTransaction.Amount))
	if err != nil {
		log.Println("ERROR money movement Capture (Transfer): ", err)
		return nil, status.Error(codes.Internal, "failed to transfer fund")
	}

	merchantWallet, err := GetWalletByID(tx, authorizedTransaction.FinalDstMerchantWalletID)
	if err != nil {
		log.Println("ERROR money movement Capture (GetWalletByID): ", err)
		return nil, status.Error(codes.Internal, "merchant wallet not found")
	}

	err = CreateTransaction(tx, authorizedTransaction.PID, srcAccount, dstMerchantAccount, authorizedTransaction.DstUserID, merchantWallet.UserID, merchantWallet.ID, int64(authorizedTransaction.Amount))
	if err != nil {
		log.Println("ERROR money movement Capture (CreateTransaction): ", err)
		return nil, status.Error(codes.Internal, "failed to create transaction")
	}

	// end tx
	err = tx.Commit()
	if err != nil {
		log.Println("ERROR money movement Capture (tx.Commit): ", err)
		return nil, status.Error(codes.Internal, "failed commiting transaction")
	}

	go publisher.SendCaptureMessage(s.publisher, authorizedTransaction.PID, authorizedTransaction.SrcUserID, authorizedTransaction.Amount)

	return &emptypb.Empty{}, nil
}
