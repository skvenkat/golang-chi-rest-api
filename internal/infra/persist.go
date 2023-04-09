package infra

import (
	"github.com/skvenkat/golang-chi-rest-api/internal/adapters/persist"
	"github.com/skvenkat/golang-chi-rest-api/internal/core/app"
	"github.com/skvenkat/golang-chi-rest-api/internal/core/di"
	"github.com/skvenkat/golang-chi-rest-api/internal/core/outport"
)

func wirePersistPorts(
	cfg *app.Config,
	cache outport.Cache,
	di *di.DI,
) func() {
	pers := persist.NewPersistence(cfg)
	addrBook := persist.NewAddrBookAdapter(
		pers,
		cache,
	)
	di.UseCases.AddrBook = addrBook
	return pers.Close
}
