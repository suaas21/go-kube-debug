package main

import (
	"context"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/godebug/app"
	debug_chi "github.com/godebug/chi"
	ctx "github.com/godebug/context"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"
	"log"
	"net/http"
)

// settlementCmd is the godebug sub command to server
var settlementCmd = &cobra.Command{
	Use:   "settlement",
	Short: "settlement to the api server",
	RunE:  settlement,
}

func init() {
	settlementCmd.PersistentFlags().StringVarP(&cfgPath, "config", "c", "", "config file path")
}

func settlement(cmd *cobra.Command, args []string) error {
	cfgApp := app.GetApp(cfgPath)
	// define service name
	cfgApp.Service = app.ServiceSettlement
	// Set the flags for the logging package to give us the filename in the logs
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	cleanup := initTracer("true", cfgApp.OTLEndpoint, cfgApp.Service)
	defer cleanup(context.Background())

	r := chi.NewMux()

	prom := ctx.NewProm(app.ServiceSettlement).Histogram(nil)
	// set middleware for chi
	r.Use(debug_chi.Use(prom))
	// define /metrics endpoint to expose prometheus metrics
	r.Mount("/metrics", promhttp.Handler())

	r.Get("/", cfgApp.Serve)

	return http.ListenAndServe(fmt.Sprintf(":%v", cfgApp.SettlementPort), r)
}
