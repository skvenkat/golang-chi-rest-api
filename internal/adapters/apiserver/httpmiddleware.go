package apiserver

import (
	"fmt"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/skvenkat/golang-chi-rest-api/internal/core/app"
	"go.uber.org/zap"
	"net/http"
	"time"
)

const chiSugaredLogFormat = `[END] "%s %s %s" from %s - %s`

func zapLoggerMiddleware(logger *zap.SugaredLogger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			wr := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			l := logger.With(zap.String("requestId", middleware.GetReqID(r.Context())))
			ctx := app.ContextWithLogger(r.Context(), l)
			tstart := time.Now()
			defer func() {
				status := wr.Status()
				statusLabel := fmt.Sprintf("%d %s", status, http.StatusText(status))
				l.With(
					zap.String("duration", time.Since(tstart).String()),
				).Infof(chiSugaredLogFormat,
					r.Method,
					r.URL.Path,
					r.Proto,
					r.RemoteAddr,
					statusLabel,
				)
			}()
			next.ServeHTTP(wr, r.WithContext(ctx))
		}
		return http.HandlerFunc(fn)
	}
}
