package chi

import (
	"github.com/go-chi/chi"
	ctx "github.com/godebug/context"
	"github.com/godebug/utils"
	"net/http"
	"strings"
	"time"
)

type context struct {
	*ctx.Context
}

func (c *context) do(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		wi := &utils.ResponseWriterInterceptor{
			StatusCode:     http.StatusOK,
			ResponseWriter: w,
		}

		chiCtx := chi.RouteContext(r.Context())

		start := time.Now()
		defer func() {
			path := strings.Join(chiCtx.RoutePatterns, "")
			for {
				nPath := strings.Replace(path, "/*/", "/", 1)
				if nPath == path {
					break
				}
				path = nPath
			}
			c.Push(path, r.Method, wi.StatusCode, start)
		}()
		next.ServeHTTP(wi, r)
	})
}

func Wrap(prom *ctx.Prom) func(http.Handler) http.Handler {
	c := &context{prom.GetContext()}
	return c.do
}

func Use(prom *ctx.Prom) func(http.Handler) http.Handler {
	c := &context{prom.GetContext()}
	return c.do
}
