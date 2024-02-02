package main

import (
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"logger/data"
	"logger/logs"
	"net"
)

type LogServer struct {
	logs.UnimplementedLogServiceServer //Ensure backwards compatibility
	Models                             data.Models
}

func (l *LogServer) WriteLog(ctx context.Context, req *logs.LogRequest) (*logs.LogResponse, error) {
	input := req.GetLogEntry()

	if input == nil {
		res := &logs.LogResponse{Result: "no input"}
		return res, errors.New("no input provided")
	}

	// write the log
	logEntry := data.LogEntry{
		Name: input.Name,
		Data: input.Data,
	}

	err := l.Models.LogEntry.Insert(logEntry)
	if err != nil {
		res := &logs.LogResponse{Result: "failed"}
		return res, err
	}

	// return a response
	res := &logs.LogResponse{Result: "logged"}
	return res, nil
}

func (app Config) gRPCListen() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", grpcPort))
	if err != nil {
		log.Fatalln("Failed to listen for gRPC:", err)
	}
	s := grpc.NewServer()

	logs.RegisterLogServiceServer(s, &LogServer{Models: app.Models})

	log.Println("gRPC Server on port", grpcPort)

	if err := s.Serve(lis); err != nil {
		log.Fatalln("Failed to listen for gRPC:", err)
	}
}
