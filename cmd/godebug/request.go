package main

import (
	"contrib.go.opencensus.io/exporter/ocagent"
	"fmt"
	"github.com/godebug/config"
	"github.com/spf13/cobra"
	"go.opencensus.io/trace"
	"io/ioutil"
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

	oce, err := ocagent.NewExporter(
		ocagent.WithInsecure(),
		ocagent.WithReconnectionPeriod(5*time.Second),
		ocagent.WithAddress(cfgApp.OCAgentHost),
		ocagent.WithServiceName("request"))
	if err != nil {
		return err
	}
	trace.RegisterExporter(oce)
	trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})

	handle("/request", func(w http.ResponseWriter, r *http.Request) {
		for i := 0; i < 100; i++ {
			svc := fmt.Sprintf("http://%s", cfgApp.Svc)
			if i%2 == 0 {
				svc = fmt.Sprintf("%s/env", svc)
			}
			err := GetRequest(svc)
			if err != nil {
				log.Fatalf(err.Error())
			}
		}
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, `request config info: %v!`, cfgApp)
	})

	return http.ListenAndServe(fmt.Sprintf(":%v", cfgApp.RequestPort), nil)
}

func GetRequest(url string) error {
	c := http.Client{Timeout: time.Duration(5) * time.Second}
	resp, err := c.Get(fmt.Sprintf("%s", url))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	log.Printf("Body : %s", body)
	return nil
}
