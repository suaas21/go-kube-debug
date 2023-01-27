package main

import (
	"contrib.go.opencensus.io/exporter/ocagent"
	"fmt"
	"github.com/go-chi/chi"
	debug_chi "github.com/godebug/chi"
	"github.com/godebug/config"
	"github.com/godebug/context"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/trace"
	"log"
	"net/http"
	"time"
)

// requestCmd is the request sub command to request api server
var requestCmd = &cobra.Command{
	Use:   "request",
	Short: "request to the api server",
	RunE:  request,
}

func init() {
	requestCmd.PersistentFlags().StringVarP(&cfgPath, "config", "c", "", "config file path")
}

func request(cmd *cobra.Command, args []string) error {
	cfgApp := config.GetApp(cfgPath)
	// Set the flags for the logging package to give us the filename in the logs
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	r := chi.NewMux()

	prom := context.NewProm("godebug").Histogram(nil)
	// set middleware for chi
	r.Use(debug_chi.Use(prom))
	// define /metrics endpoint to expose prometheus metrics
	r.Mount("/metrics", promhttp.Handler())

	oce, err := ocagent.NewExporter(
		ocagent.WithInsecure(),
		ocagent.WithReconnectionPeriod(5*time.Second),
		ocagent.WithAddress(cfgApp.OCAgentHost),
		ocagent.WithServiceName("request"))
	if err != nil {
		return err
	}
	trace.RegisterExporter(oce)

	r.Get("/request", cfgApp.Request)

	// serve http octhttp handler
	return http.ListenAndServe(fmt.Sprintf(":%v", cfgApp.RequestPort), &ochttp.Handler{
		Handler: r,
		GetStartOptions: func(r *http.Request) trace.StartOptions {
			if r.Method == http.MethodOptions || r.URL.Path == "/metrics" {
				return trace.StartOptions{
					Sampler:  trace.NeverSample(),
					SpanKind: trace.SpanKindServer,
				}
			}
			return trace.StartOptions{}
		},
	})
}
