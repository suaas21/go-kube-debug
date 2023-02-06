package main

import (
	"contrib.go.opencensus.io/exporter/ocagent"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/godebug/app"
	debug_chi "github.com/godebug/chi"
	"github.com/godebug/context"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/trace"
	"log"
	"net/http"
	"time"
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

	r := chi.NewMux()

	prom := context.NewProm(app.ServiceWarden).Histogram(nil)
	// set middleware for chi
	r.Use(debug_chi.Use(prom))
	// define /metrics endpoint to expose prometheus metrics
	r.Mount("/metrics", promhttp.Handler())

	oce, err := ocagent.NewExporter(
		ocagent.WithInsecure(),
		ocagent.WithReconnectionPeriod(5*time.Second),
		ocagent.WithAddress(cfgApp.OCAgentHost),
		ocagent.WithServiceName(app.ServiceWarden))
	if err != nil {
		return err
	}
	trace.RegisterExporter(oce)
	// if ingress is not specified for root span then undo the below line
	//trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})

	r.Get("/", cfgApp.Serve)

	return http.ListenAndServe(fmt.Sprintf(":%v", cfgApp.WardenPort), &ochttp.Handler{
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
