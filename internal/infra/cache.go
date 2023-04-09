package infra

import (
	"fmt"
	"github.com/skvenkat/golang-chi-rest-api/internal/adapters/cache"
	"github.com/skvenkat/golang-chi-rest-api/internal/core/app"
	"github.com/skvenkat/golang-chi-rest-api/internal/core/di"
	"github.com/skvenkat/golang-chi-rest-api/internal/core/outport"
)

func wireCachePorts(cfg *app.Config, _ *di.DI) (outport.Cache, func()) {
	switch cfg.Cache.Type {
	case "none":
		return cache.NewNoCache(), func() {}
	case "inmem":
		return cache.NewInMemCache(), func() {}
	default:
		panic(fmt.Sprintf("unknown cache type: %s", cfg.Cache.Type))
	}
}
