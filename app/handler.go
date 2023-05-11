package app

import (
	"bytes"
	"context"
	"fmt"
	"github.com/godebug/responses"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"time"
)

func SendRequest(ctx context.Context, method string, url string, data []byte) (*http.Response, error) {
	request, err := http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("create request error: %w", err)
	}

	client := http.Client{
		// Wrap the Transport with one that starts a span and injects the span context
		// into the outbound request headers.
		Transport: otelhttp.NewTransport(http.DefaultTransport),
		Timeout:   10 * time.Second,
	}

	// print req field to console
	reqDump, _ := httputil.DumpRequestOut(request, true)
	fmt.Printf("REQUEST:\n%s", string(reqDump))

	return client.Do(request)
}

func GetRequest(url string, r *http.Request) error {
	resp, err := SendRequest(r.Context(), "GET", url, nil)
	if err != nil {
		return err
	}

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

	//if c.Service == ServiceCore {
	//	span := trace.SpanFromContext(r.Context())
	//	span.SetAttributes()
	//
	//	ctx, span := tracer.Start(r.Context(), "update user amount")
	//	defer span.End()
	//	ctx := r.Context()
	//	time.Sleep(2 * time.Second)
	//	_, span1 := trace.StartSpan(ctx, "sleep_2_min")
	//	defer span1.End()
	//
	//	time.Sleep(2 * time.Second)
	//	_, span2 := trace.StartSpan(ctx, "sleep_4_min")
	//	defer span2.End()
	//
	//	time.Sleep(6 * time.Second)
	//	_, span3 := trace.StartSpan(ctx, "sleep_6_min")
	//	defer span3.End()
	//}

	_ = responses.ServeJSON(w, http.StatusOK, fmt.Sprintf("response from service: %v ok!", c.Service), c)
	w.WriteHeader(http.StatusOK)
}
