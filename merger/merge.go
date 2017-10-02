package merger

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
)

type Output struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
}

type Data struct {
	ResultType string   `json:"resultType"`
	Result     []Result `json:"result"`
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
func MergeNaively(output *[]byte, merges ...*[]byte) error {
	var err error
	var result Output
	var msg_type string
	var data_msg json.RawMessage
	var data Data
	var meta_data []string

	for _, a := range merges {
		if msg_type == "" {
			msg_type, err = identifyMsgType(a, data_msg)
			if err != nil && len(merges) == 1 {
				return err
			}
		}

		tmp_result := Output{Data: &data_msg}
		if err = json.Unmarshal(*a, &tmp_result); err != nil {
			if err != nil && len(merges) == 1 {
				return err
			}
		}

		if tmp_result.Status == "success" {
			result = tmp_result
		}
		switch msg_type {
		case "vector":
			var tv Data
			err = json.Unmarshal(data_msg, &tv)
			if data.Result == nil {
				data = tv
				continue
			}

			for _, val := range tv.Result {
				data.Result = append(data.Result, val)
			}
		case "matrix":
			var tm Data
			err = json.Unmarshal(data_msg, &tm)
			if data.Result == nil {
				data = tm
				continue
			}
			for _, val := range tm.Result {
				index, ok := matrixIndex(&data.Result, val.Metric)
				if ok {
					for _, v := range val.Values {
						if !isInMatrix(reflect.ValueOf(v[0]).Interface().(float64), &data.Result[index].Values) {
							data.Result[index].Values = append(data.Result[index].Values, v)
							sort.Sort(data.Result[index].Values)
						}
					}
				} else {
					data.Result = append(data.Result, val)
				}
			}
		case "meta":
			var d []string
			err = json.Unmarshal(data_msg, &d)
			if meta_data == nil {
				meta_data = d
				continue
			}
			for _, val := range d {
				meta_data = appendUnique(meta_data, val)
			}
		default:
			fmt.Println("Oops: Don't know this ResultType yet: ", msg_type)
		}
	}

	switch msg_type {
	case "vector", "matrix":
		result.Data = data
	case "meta":
		sort.Strings(meta_data)
		result.Data = meta_data
	}

	*output, err = json.Marshal(result)

	return err
}

func appendUnique(data []string, val string) []string {
	for _, v := range data {
		if v == val {
			return data
		}
	}
	return append(data, val)
}

func identifyMsgType(a *[]byte, data_msg json.RawMessage) (string, error) {
	var d Data
	rs := Output{Data: &data_msg}
	if err := json.Unmarshal(*a, &rs); err != nil {
		if err != nil {
			return "", err
		}
	}
	json.Unmarshal(data_msg, &d)
	switch d.ResultType {
	case "":
		return "meta", nil
	default:
		return d.ResultType, nil
	}
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
