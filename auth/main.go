package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/airlangga-hub/microservices-payment/auth/pb"
	"google.golang.org/grpc"
)

const (
	dbDriver = "mysql"

	dbUser     = "auth_user"
	dbPassword = "password"

	dbName = "auth"
)

func main() {
	dsn := fmt.Sprintf("%s:%s@tcp(localhost:3306)/%s", dbUser, dbPassword, dbName)

	db, err := sql.Open(dbDriver, dsn)
	if err != nil {
		log.Println(err)
		return
	}

	defer func() {
		if err := db.Close(); err != nil {
			log.Println("Error closing auth db: ", err)
		}
	}()

	if err := db.Ping(); err != nil {
		log.Println(err)
		return
	}

	// grpc server
	key := os.Getenv("SIGNING_KEY")
	s := grpc.NewServer()
	pb.RegisterAuthServiceServer(s, NewServer(db, key))

	// listen and serve
	lis, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Println("Error listening to port 7000: ", err)
		return
	}

	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

		sig := <-sigChan
		log.Printf("Signal %v received. Shutting down gRPC....\n", sig)

		s.GracefulStop()
	}()

	log.Println("Listening to port 9000.....")
	if err := s.Serve(lis); err != nil {
		log.Println("FATAL: error serving port 9000: ", err)
		log.Println("Exiting main....")
	}
}
