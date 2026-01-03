package main

import (
	"context"
	"database/sql"

	"github.com/airlangga-hub/microservices-payment/money_movement/pb"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Server struct {
	pb.UnimplementedMoneyMovementServiceServer
	db *sql.DB
}

func NewServer(db *sql.DB) *Server {
	return &Server{db: db}
}

func (s *Server) Authorize(ctx context.Context, r *pb.AuthorizeRequest) (*pb.AuthorizeResponse, error)

func (s *Server) Capture(ctx context.Context, r *pb.CaptureRequest) (*emptypb.Empty, error)
