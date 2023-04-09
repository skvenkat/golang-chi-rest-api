package internal

import (
	"github.com/go-chi/render"
	"net/http"
)

var AppVersion = &VersionRest{
	Service: "rest-chi/http",
	Version: "0.1.0",
	Build:   "1",
}

func GetVersion() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		render.Status(r, http.StatusOK)
		_ = render.Render(w, r, AppVersion)
	}
}

func (rd *VersionRest) Render(w http.ResponseWriter, r *http.Request) error {
	// Pre-processing before a response is marshalled and sent across the wire
	return nil
}
