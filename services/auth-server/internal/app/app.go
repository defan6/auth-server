package app

import (
	"log/slog"
	grpcapp "sso/internal/app/grpc"
	"sso/internal/config"
	"time"
)

type App struct {
	GRPCSrv *grpcapp.App
}

func New(
	log *slog.Logger,
	grpcPort int,
	tokenSecret []byte,
	issuer string,
	tokenTTL time.Duration,
	dbConfig *config.DBConfig,
) *App {
	grpcApp := grpcapp.New(log, grpcPort, tokenSecret, issuer, tokenTTL, dbConfig)
	return &App{
		GRPCSrv: grpcApp,
	}
}
