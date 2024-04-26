package postg_test

import (
	"context"
	"culc/iternal/config"
	"fmt"
	"net"
	"testing"

	authv1 "github.com/ragnack97/protoculc/gen/go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Postgres struct {
	T        *testing.T
	cfg      config.Config
	AuthUser authv1.AuthClient
}

func New(t *testing.T) (context.Context, *Postgres) {

	cfg := config.LoadConfig("../../config/config.yaml")
	ctx, cancelCtx := context.WithTimeout(context.Background(), cfg.GRPC.TimeOut)
	t.Cleanup(func() {
		t.Helper()
		cancelCtx()
	})
	cc, err := grpc.DialContext(context.Background(),
		net.JoinHostPort("127.0.0.1", fmt.Sprint(cfg.GRPC.Port)), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("grpc server failed connect %e", err)
	}
	return ctx, &Postgres{
		T:        t,
		cfg:      *cfg,
		AuthUser: authv1.NewAuthClient(cc),
	}
}
