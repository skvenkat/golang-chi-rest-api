package apiserver

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/skvenkat/golang-chi-rest-api/internal/adapters/apiserver/internal"
	"github.com/skvenkat/golang-chi-rest-api/internal/core/di"
	"go.uber.org/zap"
	"net/http"
)

func Start(_ context.Context, di *di.DI) {

	listenAddr := fmt.Sprintf(":%d", di.Config.Server.Port)
	zap.S().Infof("Starting http server on %s", listenAddr)

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(zapLoggerMiddleware(zap.S()))
	r.Use(middleware.Recoverer)

	apiRoutes(r, di)

	func() {
		fs := http.FileServer(http.Dir("web/dist"))
		// this is to load home page
		r.Handle("/*", http.StripPrefix("/", fs))
		// all secondary pages must be in /webapp route (or some other sub-route name, but
		// they should not be in root / route, because otherwise you might have issues
		// with React routes, when direct links to pages are not going to work.
		r.Handle("/webapp/*", rerouteToRoot(fs))
	}() // file server handler to serve web application

	srv := NewHttpServer(listenAddr, r)
	srv.ShutdownCallback = func() {
		zap.S().Info("Cleaning up resources")
		di.Close()
		zap.S().Infof("Resources has been cleaned up")
	}
	zap.S().Info("Starting HTTP server...")
	go func() {
		srv.start()
	}()

	srv.waitWithGracefulShutdown()
}

func apiRoutes(mux *chi.Mux, di *di.DI) {
	mux.Get("/api/version", internal.GetVersion())
	mux.Route("/api/contacts", func(r chi.Router) {
		r.Post("/", internal.CreateContact(di.UseCases))
		r.Get("/", internal.ListAllContacts(di.UseCases))

		r.Route("/{contactId}", func(r chi.Router) {
			r.Get("/", internal.GetContact(di.UseCases))
			r.Put("/", internal.UpdateContact(di.UseCases))
			r.Delete("/", internal.DeleteContact(di.UseCases))
		})
	})
}
