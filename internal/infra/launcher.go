package infra

import (
	"context"
	"github.com/skvenkat/golang-chi-rest-api/internal/adapters/apiserver"
	"github.com/skvenkat/golang-chi-rest-api/internal/core/app"
	"go.uber.org/zap"
)

func Start(deployment string) {
	logger, _ := zap.NewDevelopment()
	zap.ReplaceGlobals(logger)
	ctx := app.ContextWithLogger(context.Background(), zap.S())

	cfg := app.LoadConfig(deployment)
	di := wireDependencies(cfg)
	apiserver.Start(ctx, di)
}
