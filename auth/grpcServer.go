package main

import (
	"context"

	"github.com/airlangga-hub/microservices-payment/auth/pb"
)

type Server struct {
	pb.UnimplementedAuthServiceServer
}

func (s *Server) GetToken(ctx context.Context, credentials *pb.Credentials) (*pb.Token, error)

func (s *Server) ValidateToken(ctx context.Context, token *pb.Token) (*pb.User, error)
