package infra

import (
	"github.com/skvenkat/golang-chi-rest-api/internal/core/app"
	"github.com/skvenkat/golang-chi-rest-api/internal/core/di"
	"github.com/skvenkat/golang-chi-rest-api/internal/core/usecase"
	"go.uber.org/zap"
)

func wireDependencies(cfg *app.Config) *di.DI {
	zap.S().Info("Initialize DI objects")
	newDI := &di.DI{
		Config:   cfg,
		UseCases: &usecase.UseCases{},
	}

	cache, cacheCleanup := wireCachePorts(cfg, newDI)

	persistCleanup := wirePersistPorts(
		cfg,
		cache,
		newDI,
	)

	newDI.Close = func() {
		zap.S().Info("Performing cleanup of all initialized DI objects")
		persistCleanup()
		cacheCleanup()
	}
	return newDI
}
