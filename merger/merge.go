package merger

import (
	"encoding/json"
)

type PrometheusResult struct {
	Status string `json:"status"`
	Data   struct {
		ResultType string `json:"resultType"`
		Result     []struct {
			Metric interface{}   `json:"metric,omitempty"`
			Values []interface{} `json:"values,omitempty"`
			Value  interface{}   `json:"value,omitempty"`
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
				switch {
				case as.Data.ResultType == "vector":
					ts.Data.Result = append(ts.Data.Result, v)
				case as.Data.ResultType == "matrix":
					for _, v := range as.Data.Result[i].Values {
						ts.Data.Result[i].Values = append(ts.Data.Result[i].Values, v)
					}
				}
			}
		}
	}

	*result, err = json.Marshal(ts)

	return err
}