package main

import (
	"log"
	"net/http"

	authpb "github.com/airlangga-hub/microservices-payment/gateway/auth"
	mmpb "github.com/airlangga-hub/microservices-payment/gateway/money_movement"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var authClient authpb.AuthServiceClient
var mmClient mmpb.MoneyMovementServiceClient

func main() {
	authConn, err := grpc.NewClient("auth:9000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Println("FATAL: error creating authConn: ", err)
		return
	}
	defer authConn.Close()

	authClient = authpb.NewAuthServiceClient(authConn)

	mmConn, err := grpc.NewClient("money_movement:7000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Println("FATAL: error creating mmConn: ", err)
		return
	}
	defer mmConn.Close()

	mmClient = mmpb.NewMoneyMovementServiceClient(mmConn)

	http.HandleFunc("/login", login)
	http.HandleFunc("/customer/payment/authorize", customerPaymentAuthorize)
	http.HandleFunc("/customer/payment/capture", customerPaymentCapture)
}

func login(w http.ResponseWriter, r *http.Request)

func customerPaymentAuthorize(w http.ResponseWriter, r *http.Request)

func customerPaymentCapture(w http.ResponseWriter, r *http.Request)
