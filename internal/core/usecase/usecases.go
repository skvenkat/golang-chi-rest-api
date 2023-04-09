package usecase

import (
	"github.com/skvenkat/golang-chi-rest-api/internal/core/outport"
)

type UseCases struct {
	AddrBook outport.AddrBook
	// other output/secondary ports can be added here
}
