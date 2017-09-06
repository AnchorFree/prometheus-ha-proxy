package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/goware/urlx"

	"github.com/Sirupsen/logrus"
	"github.com/anchorfree/prometheus-ha-proxy/merger"
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
	Logger *logrus.Entry
	URLs   []*url.URL
}

func main() {
	addresses := os.Args[1:]
	logger := logrus.WithFields(logrus.Fields{
		"app": "prometheus-ha-proxy",
	})

	var data PromData
	data.Logger = logger
	for _, a := range addresses {
		url, err := urlx.Parse(a)
		if err != nil {
			// we could not do parse address url
			// TODO: logging
			logger.Warn("could not parse address due to", err)
			continue
		}
		data.URLs = append(data.URLs, url)
	}

	http.HandleFunc("/", data.PrometheusProxy) // set router
	err := http.ListenAndServe(":9090", nil)   // set listen port
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func (p *PromData) PrometheusProxy(w http.ResponseWriter, r *http.Request) {
	logger := p.Logger
	logger.Info(r.URL.String())
	var buffers []*[]byte
	ch := make(chan BackendOutput)

	if r.Method == "GET" {
		for _, a := range p.URLs {
			go PromGet(logger, a.ResolveReference(r.URL).String(), ch)
		}

		var cnt int
		for out := range ch {
			cnt++
			if out.Err != nil {
				continue
			}
			buffer := out.Body
			buffers = append(buffers, &buffer)
			if cnt == len(p.URLs) {
				close(ch)
			}
		}

		merged := new([]byte)
		merger.MergeNaively(merged, buffers...)
		w.Write(*merged)
	}
}

func PromGet(logger *logrus.Entry, url string, ch chan BackendOutput) {
	var out BackendOutput

	res, err := http.Get(url)
	if err != nil {
		logger.Warning("could not query ", url, "due to: ", err)
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
