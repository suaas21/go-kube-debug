package main

import (
	"contrib.go.opencensus.io/exporter/ocagent"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/godebug/config"
	"github.com/spf13/cobra"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/trace"
	"log"
	"net/http"
	"time"
)

// reqCmd is the request sub command to request api server
var reqCmd = &cobra.Command{
	Use:   "req",
	Short: "request to the api server",
	RunE:  req,
}

func init() {
	reqCmd.PersistentFlags().StringVarP(&cfgPath, "config", "c", "", "config file path")
}

func req(cmd *cobra.Command, args []string) error {
	cfgApp := config.GetApp(cfgPath)
	// Set the flags for the logging package to give us the filename in the logs
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	r := chi.NewMux()

	oce, err := ocagent.NewExporter(
		ocagent.WithInsecure(),
		ocagent.WithReconnectionPeriod(5*time.Second),
		ocagent.WithAddress(cfgApp.OCAgentHost),
		ocagent.WithServiceName("req"))
	if err != nil {
		return err
	}
	trace.RegisterExporter(oce)
	trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})

	r.Get("/req", cfgApp.Req)
	r.Get("/", cfgApp.Ok)

	// serve http octhttp handler
	return http.ListenAndServe(fmt.Sprintf(":%v", cfgApp.ReqPort), &ochttp.Handler{
		Handler: r,
	})
}
