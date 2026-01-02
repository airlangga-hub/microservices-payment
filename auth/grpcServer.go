package main

import (
	"context"
	"database/sql"

	"github.com/airlangga-hub/microservices-payment/auth/pb"
)

type Server struct {
	pb.UnimplementedAuthServiceServer
	db *sql.DB
}

func NewServer(db *sql.DB) *Server {
	return &Server{db: db}
}

func (s *Server) GetToken(ctx context.Context, credentials *pb.Credentials) (*pb.Token, error)

func (s *Server) ValidateToken(ctx context.Context, token *pb.Token) (*pb.User, error)
