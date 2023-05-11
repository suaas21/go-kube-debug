package main

import (
	"context"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/godebug/app"
	debugChi "github.com/godebug/chi"
	debugCtx "github.com/godebug/context"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"
	"log"
	"net/http"
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

	cleanup := initTracer("true", cfgApp.OTLEndpoint, cfgApp.Service)
	defer cleanup(context.Background())

	r := chi.NewMux()

	// Add custom prometheus
	prom := debugCtx.NewProm(app.ServiceSetting).Histogram(nil)
	r.Use(debugChi.Use(prom))
	r.Mount("/metrics", promhttp.Handler())

	r.Get("/", cfgApp.Serve)

	return http.ListenAndServe(fmt.Sprintf(":%v", cfgApp.SettingPort), r)
}
