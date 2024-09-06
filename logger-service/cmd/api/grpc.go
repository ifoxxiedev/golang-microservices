package main

import (
	"context"
	"fmt"
	"log"
	"log-service/cmd/data"
	"log-service/logs"
	"net"

	"google.golang.org/grpc"
)

type LogServer struct {
	logs.UnimplementedLogServiceServer
	Models data.Models
}

func (l *LogServer) WriteLog(ctx context.Context, req *logs.LogRequest) (*logs.LogResponse, error) {
	input := req.GetLogEntry()

	// Write Log
	logEntry := data.LogEntry{
		Name: input.Name,
		Data: input.Data,
	}

	err := l.Models.LogEntry.Insert(logEntry)
	if err != nil {
		return &logs.LogResponse{Result: "failed"}, err
	}

	return &logs.LogResponse{Result: "logged!"}, nil
}

func (app *Config) gRPCListen() {
	lst, err := net.Listen("tcp", fmt.Sprintf(":%s", gRpcPort))
	if err != nil {
		log.Fatalf("Failed to listen gRPC: %v", err)
	}

	s := grpc.NewServer()

	logs.RegisterLogServiceServer(s, &LogServer{Models: app.Models})
	log.Printf("gRPC Server started on port \n", gRpcPort)

	if err = s.Serve(lst); err != nil {
		log.Fatalf("Failed to start gRPC server: %v", err)
	}
}
