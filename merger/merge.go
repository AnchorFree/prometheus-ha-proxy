package merger

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
)

type PrometheusResult struct {
	Status string `json:"status"`
	Data   struct {
		ResultType string   `json:"resultType"`
		Result     []Result `json:"result"`
	} `json:"data"`
}

type Result struct {
	Metric map[string]string `json:"metric"`
	Values Values            `json:"values,omitempty"`
	Value  interface{}       `json:"value,omitempty"`
}

type Values [][]interface{}

func (v Values) Len() int {
	return len(v)
}

func (s Values) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s Values) Less(i, j int) bool {
	return reflect.ValueOf(s[i][0]).Interface().(float64) < reflect.ValueOf(s[j][0]).Interface().(float64)
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
			// https://prometheus.io/docs/querying/api/#expression-query-result-formats
			for _, val := range as.Data.Result {
				switch as.Data.ResultType {
				case "vector":
					ts.Data.Result = append(ts.Data.Result, val)
				case "matrix":
					index, ok := matrixIndex(&ts.Data.Result, val.Metric)
					if ok {
						for _, v := range val.Values {
							if !isInMatrix(reflect.ValueOf(v[0]).Interface().(float64), &ts.Data.Result[index].Values) {
								ts.Data.Result[index].Values = append(ts.Data.Result[index].Values, v)
								sort.Sort(ts.Data.Result[index].Values)
							}
						}
					} else {
						ts.Data.Result = append(ts.Data.Result, val)
					}
				default:
					fmt.Println("Oops: Don't know this ResultType yet: ", as.Data.ResultType)
				}
			}
		}
	}

	*result, err = json.Marshal(ts)

	return err
}

// index is used to store data within several runs, to seepdup the search
// get prometheus matrix results, and return true if value is present
func isInMatrix(value float64, list *Values) bool {
	result := false
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

func matrixIndex(metrics *[]Result, metric map[string]string) (int, bool) {
	for i, m := range *metrics {
		if reflect.DeepEqual(m.Metric, metric) {
			return i, true
		}
	}
	return 0, false
}
