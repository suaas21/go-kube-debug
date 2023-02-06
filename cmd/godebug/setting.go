package main

import (
	"contrib.go.opencensus.io/exporter/ocagent"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/godebug/app"
	debugChi "github.com/godebug/chi"
	debugCtx "github.com/godebug/context"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/trace"
	"log"
	"net/http"
	"time"
)

// settingCmd is the serve sub command to start the api server
var settingCmd = &cobra.Command{
	Use:   "setting",
	Short: "setting serves the api server",
	RunE:  setting,
}

func init() {
	settingCmd.PersistentFlags().StringVarP(&cfgPath, "config", "c", "", "config file path")
}

func setting(cmd *cobra.Command, args []string) error {
	cfgApp := app.GetApp(cfgPath)
	// define service name
	cfgApp.Service = app.ServiceSetting
	// Set the flags for the logging package to give us the filename in the logs
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	r := chi.NewMux()

	// Add custom prometheus
	prom := debugCtx.NewProm(app.ServiceSetting).Histogram(nil)
	r.Use(debugChi.Use(prom))
	r.Mount("/metrics", promhttp.Handler())

	// Add tracing data exporter
	oce, err := ocagent.NewExporter(
		ocagent.WithInsecure(),
		ocagent.WithReconnectionPeriod(5*time.Second),
		ocagent.WithAddress(cfgApp.OCAgentHost),
		ocagent.WithServiceName(app.ServiceSetting))
	if err != nil {
		return err
	}
	trace.RegisterExporter(oce)

	r.Get("/", cfgApp.Serve)

	return http.ListenAndServe(fmt.Sprintf(":%v", cfgApp.SettingPort), &ochttp.Handler{
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
