package di

import (
	"github.com/skvenkat/golang-chi-rest-api/internal/core/app"
	"github.com/skvenkat/golang-chi-rest-api/internal/core/usecase"
)

type DI struct {
	Close    func()
	Config   *app.Config
	UseCases *usecase.UseCases
}
