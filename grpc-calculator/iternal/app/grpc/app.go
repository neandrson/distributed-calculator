package gRPCserver

import (
	authg "culc/iternal/grpc/auth"
	"culc/iternal/services/worker"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
)

type App struct {
	gRPCServer *grpc.Server
	port       string
}

func NewApp(port string, authinter authg.Auth, countWorker int) *App {
	gPRCServer := grpc.NewServer()

	culcserv := worker.NewCalculatorServer(countWorker)
	culcserv.AddClient("1", countWorker)
	authg.Register(gPRCServer, authinter, culcserv)
	return &App{
		gRPCServer: gPRCServer,
		port:       port,
	}
}
func (s *App) Run() error {
	l, err := net.Listen("tcp", ":"+s.port)
	if err != nil {
		return fmt.Errorf("failed to start gRPC server: %w", err)
	}
	log.Print("gRPC server is starting on port :", s.port)
	go func() {
		if err := s.gRPCServer.Serve(l); err != nil {
			log.Fatalf("failed to serve gRPC server: %v", err)
		}
	}()

	return nil
}
func (a *App) Stop() {
	a.gRPCServer.GracefulStop()
}
