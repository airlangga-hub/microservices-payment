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
	db  *sql.DB
	key string
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
		WHERE
			email = ?;
		`,
	)
	if err != nil {
		log.Println("ERROR auth GetToken (db.Prepare): ", err)
		return nil, status.Error(codes.Internal, "error creating token")
	}

	var u User

	if err := stmt.QueryRow(credentials.Email).Scan(&u.Email, &u.Password); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Error(codes.Unauthenticated, "invalid credentials")
		}

		log.Println("ERROR auth GetToken (stmt.QueryRow): ", err)
		return nil, status.Error(codes.Internal, "error creating token")
	}

	jwt, err := CreateJWT(u.Email, []byte(s.key))
	if err != nil {
		log.Println("ERROR auth GetToken (CreateJWT): ", err)
		return nil, status.Error(codes.Internal, "error creating token")
	}

	return &pb.Token{Jwt: jwt}, nil
}

func (s *Server) ValidateToken(ctx context.Context, token *pb.Token) (*pb.User, error) {

	userID, err := ValidateJWT(token.Jwt, []byte(s.key))
	if err != nil {
		log.Println("ERROR auth ValidateToken (ValidateJWT): ", err)
		return nil, status.Error(codes.Unauthenticated, "invalid token")
	}

	return &pb.User{UserId: userID}, nil
}
