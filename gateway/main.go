package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/feslima/common"
	pb "github.com/feslima/common/api"
	_ "github.com/joho/godotenv/autoload"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	httpAddr         = common.EnvString("HTTP_ADDR", ":8080")
	orderServiceAddr = common.EnvString("ORDER_GRPC_ADDR", ":3000")
)

func main() {
	logHandler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelInfo})
	logger := slog.New(logHandler)

	conn, err := grpc.Dial(orderServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Error("Failed to connect to order service via gRPC")
		panic(err)
	}
	defer conn.Close()
	logger.Info(fmt.Sprintf("Connection with gRPC order service established at %s", orderServiceAddr))

	client := pb.NewOrderServiceClient(conn)

	mux := http.NewServeMux()
	handler := NewHandler(client)
	handler.registerRoutes(mux)

	logger.Info(fmt.Sprintf("Starting HTTP server at %s", httpAddr))

	if err := http.ListenAndServe(httpAddr, mux); err != nil {
		logger.Error("Failed to start http server")
		panic(err)
	}
}
