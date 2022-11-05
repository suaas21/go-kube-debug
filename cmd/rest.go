package main

import (
	"fmt"
	"github.com/godebug/config"
	"github.com/spf13/cobra"
	"log"
	"net/http"
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

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `api ok!!`)
	})

	http.HandleFunc("/env", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `environment varialbe are: %v!`, cfgApp)
	})

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", cfgApp.Port), nil))

	return nil
}
