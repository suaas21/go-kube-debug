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

// memberCmd is the godebug sub command to server
var memberCmd = &cobra.Command{
	Use:   "member",
	Short: "member serves the api server",
	RunE:  member,
}

func init() {
	memberCmd.PersistentFlags().StringVarP(&cfgPath, "config", "c", "", "config file path")
}

func member(cmd *cobra.Command, args []string) error {
	cfgApp := app.GetApp(cfgPath)
	// define service name
	cfgApp.Service = app.ServiceMember
	// Set the flags for the logging package to give us the filename in the logs
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	cleanup := initTracer("true", cfgApp.OTLEndpoint, cfgApp.Service)
	defer cleanup(context.Background())

	r := chi.NewMux()

	prom := debugCtx.NewProm(app.ServiceMember).Histogram(nil)
	r.Use(debugChi.Use(prom))
	r.Mount("/metrics", promhttp.Handler())

	r.Get("/", cfgApp.Serve)

	return http.ListenAndServe(fmt.Sprintf(":%v", cfgApp.MemberPort), r)
}
