package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"

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
		log.Fatalln(err)
	}

	defer func() {
		if err := db.Close(); err != nil {
			log.Println("Error closing auth db: ", err)
		}
	}()

	if err := db.Ping(); err != nil {
		log.Fatalln(err)
	}

	// grpc server
	key := os.Getenv("SIGNING_KEY")
	s := grpc.NewServer()
	pb.RegisterAuthServiceServer(s, NewServer(db, key))

	// listen and serve
	lis, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatalln("Error listening to port 9000: ", err)
	}

	log.Println("Listening to port 9000.....")
	log.Fatalln("Program terminated: ", s.Serve(lis))
}
