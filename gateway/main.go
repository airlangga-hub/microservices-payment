package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"

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

func login(w http.ResponseWriter, r *http.Request) {

	username, password, ok := r.BasicAuth()
	if !ok {
		http.Error(w, "invalid user credentials", http.StatusUnauthorized)
		return
	}

	token, err := authClient.GetToken(
		context.Background(),
		&authpb.Credentials{
			Email:    username,
			Password: password,
		},
	)
	if err != nil {
		log.Println("ERROR gateway login (GetToken): ", err)
		http.Error(w, "error creating token", http.StatusInternalServerError)
		return
	}

	_, err = w.Write([]byte(token.Jwt))
	if err != nil {
		log.Println("ERROR gateway login (w.Write): ", err)
	}
}

func customerPaymentAuthorize(w http.ResponseWriter, r *http.Request) {

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "empty auth header", http.StatusUnauthorized)
		return
	}
	if !strings.HasPrefix(authHeader, "Bearer ") {
		http.Error(w, "malformed auth header", http.StatusUnauthorized)
		return
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")

	ctx := context.Background()

	user, err := authClient.ValidateToken(
		ctx,
		&authpb.Token{Jwt: token},
	)
	if err != nil {
		http.Error(w, "invalid token", http.StatusUnauthorized)
		return
	}

	var req AuthorizeRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	res, err := mmClient.Authorize(
		ctx,
		&mmpb.AuthorizeRequest{
			CustomerWalletUserId: user.Email,
			MerchantWalletUserId: req.MerchantWalletUserId,
			Cents:                req.Cents,
			Currency:             req.Currency,
		},
	)
	if err != nil {
		log.Println("ERROR gateway customerPaymentAuthorize (mmClient.Authorize): ", err)
		http.Error(w, "transaction authorization failed", http.StatusInternalServerError)
		return
	}

	response := AuthorizeResponse{PID: res.Pid}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Println("ERROR gateway customerPaymentAuthorize (.Encode): ", err)
	}
}

func customerPaymentCapture(w http.ResponseWriter, r *http.Request)
