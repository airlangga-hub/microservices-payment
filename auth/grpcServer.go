package main

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/airlangga-hub/microservices-payment/auth/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	pb.UnimplementedAuthServiceServer
	db *sql.DB
}

func NewServer(db *sql.DB) *Server {
	return &Server{db: db}
}

func (s *Server) GetToken(ctx context.Context, credentials *pb.Credentials) (*pb.Token, error) {

	stmt, err := s.db.Prepare(
		`
		SELECT
			email,
			password
		FROM users
		WHERE email = ?;
		`,
	)
	if err != nil {
		log.Println("ERROR auth GetToken: ", err)
		return nil, status.Error(codes.Internal, "error creating token")
	}

	var u User

	if err := stmt.QueryRow(credentials.Email).Scan(&u.Email, &u.Password); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Error(codes.Unauthenticated, "invalid credentials")
		}
	}

	return &pb.Token{}, nil
}

func (s *Server) ValidateToken(ctx context.Context, token *pb.Token) (*pb.User, error)
