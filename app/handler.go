package app

import (
	"fmt"
	"github.com/godebug/responses"
	"go.opencensus.io/plugin/ochttp/propagation/b3"
	"go.opencensus.io/trace"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"time"
)

func ForwardTracingHeaderMetaData(incomingReq, outgoingReq *http.Request) {
	f := b3.HTTPFormat{}
	sc, ok := f.SpanContextFromRequest(incomingReq)
	if ok {
		f.SpanContextToRequest(sc, outgoingReq)
	}
}

func GetRequest(url string, r *http.Request) error {
	c := &http.Client{
		Timeout: time.Duration(60) * time.Second,
	}

	req, _ := http.NewRequest("GET", fmt.Sprintf("%s", url), nil)

	// forward tracing header metadata to request
	ForwardTracingHeaderMetaData(r, req)

	// print req field to console
	reqDump, _ := httputil.DumpRequestOut(req, true)
	fmt.Printf("REQUEST:\n%s", string(reqDump))

	resp, err := c.Do(req)
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

func (c *Application) Serve(w http.ResponseWriter, r *http.Request) {

	url := "http://"
	switch c.Service {
	case ServiceWarden:
		url = url + c.MemberSVC
	case ServiceMember:
		url = url + c.SettingSVC
	case ServiceSetting:
		url = url + c.SettlementSVC
	case ServiceSettlement:
		url = url + c.CoreSVC
	}

	if c.Service != ServiceCore {
		err := GetRequest(url, r)
		if err != nil {
			_ = responses.ServeJSON(w, http.StatusInternalServerError, err.Error(), nil)
			return
		}
	}

	if c.Service == ServiceCore {
		ctx := r.Context()
		time.Sleep(2 * time.Second)
		_, span1 := trace.StartSpan(ctx, "sleep_2_min")
		defer span1.End()

		time.Sleep(2 * time.Second)
		_, span2 := trace.StartSpan(ctx, "sleep_4_min")
		defer span2.End()

		time.Sleep(6 * time.Second)
		_, span3 := trace.StartSpan(ctx, "sleep_6_min")
		defer span3.End()
	}

	_ = responses.ServeJSON(w, http.StatusOK, fmt.Sprintf("response from service: %v ok!", c.Service), c)
	w.WriteHeader(http.StatusOK)
}
