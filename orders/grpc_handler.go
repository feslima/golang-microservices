package main

import (
	"context"
	"fmt"
	"log/slog"

	pb "github.com/feslima/common/api"
	"google.golang.org/grpc"
)

type grpcHandler struct {
	pb.UnimplementedOrderServiceServer
	logger *slog.Logger
}

func NewGRPCHandler(grpcServer *grpc.Server, logger *slog.Logger) *grpcHandler {
	handler := &grpcHandler{pb.UnimplementedOrderServiceServer{}, logger}
	pb.RegisterOrderServiceServer(grpcServer, handler)

	return handler
}

func (h *grpcHandler) CreateOrder(ctx context.Context, request *pb.CreateOrderRequest) (*pb.Order, error) {
	h.logger.Info(fmt.Sprintf("New order received! Order: %v", request))
	order := &pb.Order{
		ID: "1",
	}
	return order, nil
}
