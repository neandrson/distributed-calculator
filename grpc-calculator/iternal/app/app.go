package app

import (
	"context"
	gRPCServer "culc/iternal/app/grpc"
	"log"
	"net/http"

	"culc/iternal/postg"
	"culc/iternal/services/auth"
	"time"

	culcv1 "github.com/ragnack97/protoculc/gen/go"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type App struct {
	GRPCS *gRPCServer.App
}

func New(port string, user string, password string, host string, portdb string, dbname string, tokentl time.Duration, countWorker int) *App {
	storage := postg.ConnectDB(user, password, host, portdb, dbname)

	authServer := auth.New(storage, tokentl)

	grps := gRPCServer.NewApp(port, authServer, countWorker)

	return &App{
		GRPCS: grps,
	}
}

func StartHTTP(porthttp string, portgrpc string, ctx context.Context) error {
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	err := culcv1.RegisterAuthHandlerFromEndpoint(ctx, mux, "localhost:"+portgrpc, opts)
	if err != nil {
		return err
	}

	go func() {
		server := &http.Server{
			Addr: ":" + porthttp,
			Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Access-Control-Allow-Origin", "*")
				if r.Method == http.MethodOptions {
					// Устанавливаем заголовки CORS

					w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
					w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
					w.WriteHeader(http.StatusOK)
					return
				}

				// Обработка других запросов
				mux.ServeHTTP(w, r)
			}),
		}

		if err := server.ListenAndServe(); err != nil {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()

	return nil
}
