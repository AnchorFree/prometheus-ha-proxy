package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/goware/urlx"

	"github.com/anchorfree/prometheus-ha-proxy/merger"
)

type BackendOutput struct {
	Body []byte
	Err  error
}

func main() {

	http.HandleFunc("/", PrometheusProxy)    // set router
	err := http.ListenAndServe(":9090", nil) // set listen port
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func PrometheusProxy(w http.ResponseWriter, r *http.Request) {
	addresses := os.Args[1:]
	// TODO: logging
	fmt.Println(r.URL.String())
	var buffers []*[]byte
	ch := make(chan BackendOutput)

	if r.Method == "GET" {
		for _, a := range addresses {
			base, err := urlx.Parse(a)
			if err != nil {
				// we could not do parse address url
				// TODO: logging
				fmt.Println(err)
				continue
			}
			go PromGet(base.ResolveReference(r.URL).String(), ch)
		}

		var cnt int
		for out := range ch {
			cnt++
			if out.Err != nil {
				continue
			}
			buffer := out.Body
			buffers = append(buffers, &buffer)
			if cnt == len(addresses) {
				close(ch)
			}
		}

		merged := new([]byte)
		merger.MergeNaively(merged, buffers...)
		w.Write(*merged)
	}
}

func PromGet(url string, ch chan BackendOutput) {
	var out BackendOutput

	res, err := http.Get(url)
	if err != nil {
		// TODO: logging
		fmt.Println("could not query due to: ", err)
		out.Err = err
		ch <- out
		return
	}
	defer res.Body.Close()

	if res.StatusCode >= 200 && res.StatusCode < 300 {
		out.Body, out.Err = ioutil.ReadAll(res.Body)
		ch <- out
	} else {
		// Error based on status code received
		fmt.Println("Did not work out due to: ", res.StatusCode)
	}
	return
}
