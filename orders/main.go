package main

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"

	"github.com/feslima/common"
	"google.golang.org/grpc"
)

var (
	grpcAddr = common.EnvString("ORDER_GRPC_ADDR", "localhost:3000")
)

func main() {
	logHandler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelInfo})
	logger := slog.New(logHandler)

	grpcServer := grpc.NewServer()
	listener, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		logger.Error(err.Error())
	}
	defer listener.Close()

	store := NewStore()
	svc := NewService(store)
	NewGRPCHandler(grpcServer, logger)

	err = svc.CreateOrder(context.Background())
	if err != nil {
		logger.Error(err.Error())
	}
	logger.Info(fmt.Sprintf("gPRC Server started at %s", grpcAddr))

	if err := grpcServer.Serve(listener); err != nil {
		logger.Error(err.Error())
	}
}
