package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/goware/urlx"

	"github.com/AnchorFree/prometheus-ha-proxy/merger"
	"github.com/sirupsen/logrus"
)

func init() {
	// Log as JSON instead of the default ASCII formatter.
	logrus.SetFormatter(&logrus.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	logrus.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	logrus.SetLevel(logrus.InfoLevel)
}

type BackendOutput struct {
	Body []byte
	Err  error
}

type PromData struct {
	Logger    *logrus.Entry
	EndPoints []Endpoint
}

type Endpoint struct {
	URL     *url.URL
	Active  bool
	Changed time.Time
}

func (ep *Endpoint) Set(url *url.URL) bool {
	ep.URL = url
	ep.Active = true
	ep.Changed = time.Now()
	return true
}

func (ep *Endpoint) Disable() bool {
	ep.Active = false
	ep.Changed = time.Now()
	return true
}

func (ep *Endpoint) Enable() bool {
	ep.Active = true
	ep.Changed = time.Now()
	return true
}

func main() {
	addresses := os.Args[1:]
	logger := logrus.WithFields(logrus.Fields{
		"app": "prometheus-ha-proxy",
	})

	var data PromData
	var ep Endpoint
	data.Logger = logger
	for _, a := range addresses {
		logger.Debug("Parsing: ", a)
		my_url, err := urlx.Parse(a)
		if err != nil {
			// we could not do parse address url
			logger.Warn("could not parse address due to", err)
			continue
		}
		ep.Set(my_url)
		data.EndPoints = append(data.EndPoints, ep)
	}

	go data.EndpointsProbe()

	http.HandleFunc("/", data.PrometheusProxy) // set router
	http.HandleFunc("/health", HealthcheckHandler)
	err := http.ListenAndServe(":9090", nil) // set listen port
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func (p *PromData) EndpointsProbe() {
	// we don't want to probe urls
	logger := p.Logger
	t, _ := time.ParseDuration("5m")
	ticker := time.NewTicker(t)
	m, _ := url.Parse("/metrics")
	for range ticker.C {
		// every t period, do
		for i, ep := range p.EndPoints {
			//url := ep.URL.ResolveReference(m).String()
			url := ep.URL.ResolveReference(m)
			_, err := http.Get(url.String())
			switch {
			case err != nil:
				logger.Warning(ep.URL.String(), " is DOWN")
				p.EndPoints[i].Disable()
			case !ep.Active:
				p.EndPoints[i].Enable()
				logger.Debug("Enabled ", ep.URL.String())
			}
		}
	}
}

func (p *PromData) PrometheusProxy(w http.ResponseWriter, r *http.Request) {
	logger := p.Logger
	logger.Info(r.URL.String())
	var buffers []*[]byte
	ch := make(chan BackendOutput)

	if r.Method == "GET" {
		var reqs int
		for i, ep := range p.EndPoints {
			if ep.Active {
				logger.Debug("endpoint ", ep.URL.String(), " is active")
				go PromGet(logger, &p.EndPoints[i], r.URL, ch)
				reqs++
			}
		}

		var cnt int
		for out := range ch {
			cnt++
			if out.Err == nil {
				buffer := out.Body
				buffers = append(buffers, &buffer)
			}
			if cnt == reqs {
				close(ch)
			}
		}

		merged := new([]byte)
		// who needs error logging anyways, right?
		_ = merger.MergeNaively(merged, buffers...)
		_, _ = w.Write(*merged)
	}
}

func PromGet(logger *logrus.Entry, ep *Endpoint, r *url.URL, ch chan BackendOutput) {
	var out BackendOutput

	url := ep.URL.ResolveReference(r)

	logger.Debug(url)
	res, err := http.Get(url.String())
	if err != nil {
		logger.Warning("could not query ", url, "due to: ", err)
		ep.Disable()
		out.Err = err
		ch <- out
		return
	}
	defer res.Body.Close()

	if res.StatusCode >= 200 && res.StatusCode < 300 {
		out.Body, out.Err = ioutil.ReadAll(res.Body)
		ch <- out
	} else {
		logger.Warning("Result code is not success: ", res.StatusCode)
	}
	return
}

// HealthcheckHandler /health
// just returns 200 '{ "ok": true }'
func HealthcheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if _, err := fmt.Fprintf(w, "{ \"ok\": true }"); err != nil {
		fmt.Printf("Error in response writing: %#v", err)
	}
}
