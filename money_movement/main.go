package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"

	// "github.com/airlangga-hub/microservices-payment/money_movement/pb"
	"github.com/IBM/sarama"
	"github.com/airlangga-hub/microservices-payment/money_movement/pb"
	"google.golang.org/grpc"
)

const (
	dbDriver = "mysql"

	dbUser     = "money_movement_user"
	dbPassword = "password"

	dbName = "money_movement"
)

func main() {
	dsn := fmt.Sprintf("%s:%s@tcp(localhost:3306)/%s", dbUser, dbPassword, dbName)

	db, err := sql.Open(dbDriver, dsn)
	if err != nil {
		log.Fatalln(err)
	}

	defer func() {
		if err := db.Close(); err != nil {
			log.Println("Error closing db: ", err)
		}
	}()

	if err := db.Ping(); err != nil {
		log.Fatalln(err)
	}

	// create publisher
	publisher, err := sarama.NewSyncProducer([]string{"localhost:9092"}, sarama.NewConfig())
	if err != nil {
		log.Fatalln(err)
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
		log.Fatalln("Error listening to port 7000: ", err)
	}

	log.Println("Listening to port 7000.....")
	log.Fatalln("Program terminated: ", s.Serve(lis))
}
