package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	// "github.com/airlangga-hub/microservices-payment/money_movement/pb"
	"github.com/IBM/sarama"
	"github.com/airlangga-hub/microservices-payment/money_movement/pb"
	"google.golang.org/grpc"
)

const (
	dbDriver = "mysql"

	dbUser     = "root"
	dbPassword = "password"

	dbName = "money_movement"
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
			log.Println("Error closing money_movement db: ", err)
		}
	}()

	if err := db.Ping(); err != nil {
		log.Println(err)
		return
	}

	// create publisher
	publisher, err := sarama.NewSyncProducer([]string{"localhost:9092"}, sarama.NewConfig())
	if err != nil {
		log.Println(err)
		return
	}

	defer func() {
		if err := publisher.Close(); err != nil {
			log.Println("Error closing publisher: ", err)
		}
	}()

	// grpc server
	s := grpc.NewServer()
	pb.RegisterMoneyMovementServiceServer(s, NewServer(db, publisher))

	// listen and serve
	lis, err := net.Listen("tcp", ":7000")
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

	log.Println("Listening to port 7000.....")
	if err := s.Serve(lis); err != nil {
		log.Println("FATAL: error serving port 7000: ", err)
		log.Println("Exiting main....")
	}
}
