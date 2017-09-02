package merger

import (
	"encoding/json"
	"reflect"
)

type PrometheusResult struct {
	Status string `json:"status"`
	Data   struct {
		ResultType string `json:"resultType"`
		Result     []struct {
			Metric interface{}     `json:"metric,omitempty"`
			Values [][]interface{} `json:"values,omitempty"`
			Value  interface{}     `json:"value,omitempty"`
		} `json:"result"`
	} `json:"data"`
}

// takes arbitratry amount of prometheus responses in json format
// merges all of the results into one "success"-ful
func MergeNaively(result *[]byte, merges ...*[]byte) error {
	var err error
	var ts PrometheusResult

	for _, a := range merges {
		var as PrometheusResult
		if err = json.Unmarshal(*a, &as); err != nil {
			if err != nil && len(merges) == 1 {
				return err
			}
		}
		// we need one result to append values to, but
		// we don't want to increase amount of metrics
		if ts.Status == "" && as.Status == "success" {
			ts = as
			continue
		}

		if as.Status == "success" {
			for i, v := range as.Data.Result {
				// https://prometheus.io/docs/querying/api/#expression-query-result-formats
				switch {
				case as.Data.ResultType == "vector":
					ts.Data.Result = append(ts.Data.Result, v)
				case as.Data.ResultType == "matrix":
					for _, v := range as.Data.Result[i].Values {
						if !isInMatrix(reflect.ValueOf(v[0]).Interface().(float64), &ts.Data.Result[i].Values) {
							ts.Data.Result[i].Values = append(ts.Data.Result[i].Values, v)
						}
					}
				}
			}
		}
	}

	*result, err = json.Marshal(ts)

	return err
}

// index is used to store data within several runs, to seepdup the search
// get prometheus matrix results, and return true if value is present
func isInMatrix(value float64, list *[][]interface{}) bool {
	var result bool
	for _, v := range *list {
		switch reflect.TypeOf(v[0]).Name() {
		case "int":
			result = reflect.ValueOf(v[0]).Interface().(int) == int(value)
		case "float64":
			result = reflect.ValueOf(v[0]).Interface().(float64) == value
		}

		if result {
			return true
		}
	}
	return result
}
