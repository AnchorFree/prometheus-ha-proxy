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

func main() {

	http.HandleFunc("/", PrometheusProxy)    // set router
	err := http.ListenAndServe(":9090", nil) // set listen port
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func PrometheusProxy(w http.ResponseWriter, r *http.Request) {
	addresses := os.Args[1:]
	fmt.Println(r.URL.String())
	var buffers []*[]byte

	if r.Method == "GET" {
		for _, a := range addresses {
			base, err := urlx.Parse(a)
			if err != nil {
				// we could not do parse address url
				// TODO: logging
				fmt.Println(err)
				continue
			}
			res, err := http.Get(base.ResolveReference(r.URL).String())
			if err != nil {
				// TODO: logging
				fmt.Println("could not query due to: ", err)
				continue
			}

			defer res.Body.Close()
			if res.StatusCode >= 200 && res.StatusCode < 300 {
				t, err := ioutil.ReadAll(res.Body)
				if err != nil {
					continue
				}

				buffers = append(buffers, &t)
			} else {
				// Error based on status code received
				fmt.Println("Did not work out due to: ", res.StatusCode)
			}
		}

		merged := new([]byte)
		merger.MergeNaively(merged, buffers...)
		w.Write(*merged)
	}
}
