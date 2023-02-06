package app

import (
	"context"
	"go.opencensus.io/plugin/ochttp/propagation/b3"
	"net/http"
)

func Tracer(cfg *Application) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			f := b3.HTTPFormat{}
			sc, ok := f.SpanContextFromRequest(r)
			if !ok {
				next.ServeHTTP(w, r)
				return
			}
			ctx := context.Background()
			ctx = context.WithValue(ctx, f, sc)
			r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}
