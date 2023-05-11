package main

import (
	"context"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/godebug/app"
	debug_chi "github.com/godebug/chi"
	ctx "github.com/godebug/context"
	"github.com/godebug/responses"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	_ "go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"log"
	"net/http"
)

// wardenCmd is the godebug sub command to server
var wardenCmd = &cobra.Command{
	Use:   "warden",
	Short: "warden to the api server",
	RunE:  warden,
}

func init() {
	wardenCmd.PersistentFlags().StringVarP(&cfgPath, "config", "c", "", "config file path")
}

func warden(cmd *cobra.Command, args []string) error {
	cfgApp := app.GetApp(cfgPath)
	// define service name
	cfgApp.Service = app.ServiceWarden
	// Set the flags for the logging package to give us the filename in the logs
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	cleanup := initTracer("true", cfgApp.OTLEndpoint, cfgApp.Service)
	defer cleanup(context.Background())

	r := chi.NewMux()

	prom := ctx.NewProm(app.ServiceWarden).Histogram(nil)
	// set middleware for chi
	r.Use(debug_chi.Use(prom))
	// define /metrics endpoint to expose prometheus metrics
	r.Mount("/metrics", promhttp.Handler())

	r.Get("/test", func(w http.ResponseWriter, r *http.Request) {
		_ = responses.ServeJSON(w, http.StatusOK, "ok", "")
	})

	r.Get("/", cfgApp.Serve)

	return http.ListenAndServe(fmt.Sprintf(":%v", cfgApp.WardenPort), otelhttp.NewHandler(r, ""))
}
