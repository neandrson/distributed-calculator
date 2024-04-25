package main

import (
	"context"
	"culc/iternal/app"
	"culc/iternal/config"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := config.LoadConfig("./config/config.yaml")

	MainApp := app.New(cfg.GRPC.Port, cfg.Postgres.User, cfg.Postgres.Password, cfg.Postgres.Host, cfg.Postgres.Port, cfg.Postgres.DBname, cfg.Token_time, cfg.CountWorker)

	err := MainApp.GRPCS.Run()
	if err != nil {
		panic("gRPC wasnt built")
	}

	app.StartHTTP(cfg.GRPCWeb, cfg.GRPC.Port, context.Background())

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	<-stop
	log.Print("grpc stoped")
	MainApp.GRPCS.Stop()
}
