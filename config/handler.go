package config

import (
	"fmt"
	"github.com/godebug/responses"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func (c *Application) Ok(w http.ResponseWriter, r *http.Request) {
	_ = responses.ServeJSON(w, http.StatusOK, `api ok!!`, nil)
}

func (c *Application) Env(w http.ResponseWriter, r *http.Request) {
	_ = responses.ServeJSON(w, http.StatusOK, `environment variable`, c)
}

func (c *Application) Request(w http.ResponseWriter, r *http.Request) {
	fmt.Println("..........request loop started....................")
	for i := 0; i < 100; i++ {
		svc := fmt.Sprintf("http://%s", c.Svc)
		if i%2 == 0 {
			svc = fmt.Sprintf("%s/env", svc)
		}
		err := GetRequest(svc)
		if err != nil {
			log.Fatalf(err.Error())
		}
	}
	w.WriteHeader(http.StatusOK)
	_, _ = fmt.Fprintf(w, `request config info: %v!`, c)
}

func GetRequest(url string) error {
	fmt.Println("----------request individual started.............")
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

func (c *Application) Req(w http.ResponseWriter, r *http.Request) {
	fmt.Println("..................req..........................")
	svc := fmt.Sprintf("http://%s/request", c.ReqSvc)
	err := GetRequest(svc)
	if err != nil {
		log.Fatalf(err.Error())
	}

	w.WriteHeader(http.StatusOK)
	_, _ = fmt.Fprintf(w, `request config info: %v!`, c)
}
