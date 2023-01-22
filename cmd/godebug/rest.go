package main

import (
	"contrib.go.opencensus.io/exporter/ocagent"
	"fmt"
	"github.com/go-chi/chi"
	debugChi "github.com/godebug/chi"
	"github.com/godebug/config"
	debugCtx "github.com/godebug/context"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/trace"
	"log"
	"net/http"
	"time"
)

// srvCmd is the serve sub command to start the api server
var srvCmd = &cobra.Command{
	Use:   "serve",
	Short: "serve serves the api server",
	RunE:  serve,
}

func init() {
	srvCmd.PersistentFlags().StringVarP(&cfgPath, "config", "c", "", "config file path")
}

func serve(cmd *cobra.Command, args []string) error {
	cfgApp := config.GetApp(cfgPath)
	// Set the flags for the logging package to give us the filename in the logs
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	r := chi.NewMux()

	prom := debugCtx.NewProm("godebug").Histogram(nil)
	r.Use(debugChi.Use(prom))
	r.Mount("/metrics", promhttp.Handler())

	oce, err := ocagent.NewExporter(
		ocagent.WithInsecure(),
		ocagent.WithReconnectionPeriod(5*time.Second),
		ocagent.WithAddress(cfgApp.OCAgentHost),
		ocagent.WithServiceName("godebug"))
	if err != nil {
		return err
	}
	trace.RegisterExporter(oce)
	//trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})

	r.Get("/env", cfgApp.Env)
	r.Get("/", cfgApp.Ok)

	return http.ListenAndServe(fmt.Sprintf(":%v", cfgApp.Port), &ochttp.Handler{
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
