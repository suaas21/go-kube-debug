package main

import (
	"contrib.go.opencensus.io/exporter/ocagent"
	"fmt"
	"github.com/godebug/config"
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

	oce, err := ocagent.NewExporter(
		ocagent.WithInsecure(),
		ocagent.WithReconnectionPeriod(5*time.Second),
		ocagent.WithAddress(cfgApp.OCAgentHost),
		ocagent.WithServiceName("godebug"))
	if err != nil {
		return err
	}
	trace.RegisterExporter(oce)
	trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})

	handle("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `api ok!!`)
	})

	handle("/env", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `environment varialbe are: %v!`, cfgApp)
	})

	return http.ListenAndServe(fmt.Sprintf(":%v", cfgApp.Port), nil)
}

func handle(path string, h func(w http.ResponseWriter, r *http.Request)) {
	http.Handle(path, &ochttp.Handler{
		Handler: http.HandlerFunc(h),
	})
}
