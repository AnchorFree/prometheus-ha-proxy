package merger

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestVectorValuesVariadic(t *testing.T) {
	a := []byte(`{"status":"success","data":{"resultType":"vector","result":[{"metric":{"host_name":"hss4"},"value":[1504157313.787,"24"]},{"metric":{"host_name":"hss508"},"value":[1504157313.787,"32"]}]}}`)
	b := []byte(`{"status":"success","data":{"resultType":"vector","result":[{"metric":{"host_name":"hss4"},"value":[1504157435.713,"20"]},{"metric":{"host_name":"bhss-nyc2-prod-80"},"value":[1504157435.713,"2"]}]}}`)
	c := []byte(`{"status":"success","data":{"resultType":"vector","result":[{"metric":{"host_name":"hss4"},"value":[1504157435.713,"20"]},{"metric":{"host_name":"bhss-nyc2-prod-80"},"value":[1504157435.713,"2"]}]}}`)
	result := []byte(`{"status":"success","data":{"resultType":"vector","result":[{"metric":{"host_name":"hss4"},"value":[1504157313.787,"24"]},{"metric":{"host_name":"hss508"},"value":[1504157313.787,"32"]},{"metric":{"host_name":"hss4"},"value":[1504157435.713,"20"]},{"metric":{"host_name":"bhss-nyc2-prod-80"},"value":[1504157435.713,"2"]}, {"metric":{"host_name":"hss4"},"value":[1504157435.713,"20"]}, {"metric":{"host_name":"bhss-nyc2-prod-80"},"value":[1504157435.713,"2"]}]}}`)

	var o1 interface{}
	var o2 interface{}
	var err error
	z := new([]byte)

	err = MergeNaively(z, &a, &b, &c)
	if err != nil {
		t.Error("mergeJson failed with", err)
	}

	json.Unmarshal(*z, &o1)
	json.Unmarshal(result, &o2)

	if !reflect.DeepEqual(o1, o2) {
		t.Error("Expected: ", string(result), "got", string(*z))
	}
}

func TestMatrixValuesVariadic(t *testing.T) {
	a := []byte(`{
	"status": "success",
	"data": {
		"resultType": "matrix",
		"result": [{
			"metric": {},
			"values": [
				[1504286220, "0.11923089739676468"],
				[1504286224, "0.11923089739676468"]
			]
		}]
	}
}`)
	b := []byte(`{
	"status": "success",
	"data": {
		"resultType": "matrix",
		"result": [{
			"metric": {},
			"values": [
				[1504286220, "0.12527084236321268"],
				[1504286224, "0.12527084236321268"],
				[1504286228, "0.12527084236321269"]
			]
		}]
	}
}`)
	result := []byte(`{
	"status": "success",
	"data": {
		"resultType": "matrix",
		"result": [{
			"metric": {},
			"values": [
				[1504286220, "0.11923089739676468"],
				[1504286224, "0.11923089739676468"],
				[1504286228, "0.12527084236321269"]
			]
		}]
	}
}`)

	var o1 interface{}
	var o2 interface{}
	var err error
	z := new([]byte)

	err = MergeNaively(z, &a, &b)
	if err != nil {
		t.Error("mergeJson failed with", err)
	}

	json.Unmarshal(*z, &o1)
	json.Unmarshal(result, &o2)

	if !reflect.DeepEqual(o1, o2) {
		t.Error("Expected: ", string(result), "got", string(*z))
	}
}

func TestMatrixEmpty(t *testing.T) {
	a := []byte(`{
		"status": "success",
		"data": {
			"resultType": "matrix",
			"result": []
		}
	}`)

	b := []byte(`{
		"status": "success",
		"data": {
			"resultType": "matrix",
			"result": []
		}
	}`)
	result := []byte(`{
		"status": "success",
		"data": {
			"resultType": "matrix",
			"result": []
		}
	}`)

	var o1 interface{}
	var o2 interface{}
	var err error
	z := new([]byte)

	err = MergeNaively(z, &a, &b)
	if err != nil {
		t.Error("mergeJson failed with", err)
	}

	json.Unmarshal(*z, &o1)
	json.Unmarshal(result, &o2)

	if !reflect.DeepEqual(o1, o2) {
		t.Error("Expected: ", string(result), "got", string(*z))
	}
}

func TestIsInMatrix(t *testing.T) {
	i := Values{[]interface{}{1504289180.23, "0.13633787792895385"}}
	result := isInMatrix(1504289180.23, &i)
	if !result {
		t.Error("Expected: True got", result)
	}
}

func TestInvalidJSON(t *testing.T) {
	a := []byte(`{
	"status": "success",
	"data": {
		"resultType": "matrix",
		"result": [{
			"metric": {},
			"values": [
				[1504286220, "0.11923089739676468"],
				[1504286224, "0.11923089739676468"]
			]
		}]
	}
}`)
	b := []byte(`{
	"status": "success",
	"data": {
		"resultType": "matrix",
		"result": [{
			"metric": {},
			"values": [
				[1504286220, "0.12527084236321268"],
				[1504286224, "0.12527084236321268"]
				[1504286228, "0.12527084236321269"]
			]
		}]
	}
}`)
	result := []byte(`{
	"status": "success",
	"data": {
		"resultType": "matrix",
		"result": [{
			"metric": {},
			"values": [
				[1504286220, "0.11923089739676468"],
				[1504286224, "0.11923089739676468"]
			]
		}]
	}
}`)

	var o1 interface{}
	var o2 interface{}
	var err error
	z := new([]byte)

	err = MergeNaively(z, &a, &b)
	if err != nil {
		t.Error("mergeJson failed with", err)
	}

	json.Unmarshal(*z, &o1)
	json.Unmarshal(result, &o2)

	if !reflect.DeepEqual(o1, o2) {
		t.Error("Expected: ", string(result), "got", string(*z))
	}
}

func TestVariableSetOfMetrics(t *testing.T) {
	a := []byte(`
	{
	"status": "success",
	"data": {
		"resultType": "matrix",
		"result": [{
			"metric": {
				"host_name": "lhss-nyc-prod-1"
			},
			"values": [
				[1504283307, "5.73"],
				[1504360587, "2.56"]
			]
		}, {
			"metric": {
				"host_name": "lhss-nyc-prod-4"
			},
			"values": [
				[1504293387, "7.2"],
				[1504313547, "2.33"],
				[1504358187, "3.44"]
			]
		}]
	}
}`)

	b := []byte(`{
	"status": "success",
	"data": {
		"resultType": "matrix",
		"result": [{
			"metric": {
				"host_name": "lhss-nyc-prod-39"
			},
			"values": [
				[1504292907, "4.94"],
				[1504357707, "3.22"]
			]
		}, {
			"metric": {
				"host_name": "lhss-nyc-prod-4"
			},
			"values": [
				[1504293387, "7.2"],
				[1504323147, "2.81"],
				[1504335627, "3.54"]
			]
		}, {
			"metric": {
				"host_name": "lhss-nyc-prod-40"
			},
			"values": [
				[1504280427, "5.26"],
				[1504293867, "4.89"]
			]
		}]
	}
}`)

	result := []byte(`
	{
	"status": "success",
	"data": {
		"resultType": "matrix",
		"result": [{
			"metric": {
				"host_name": "lhss-nyc-prod-1"
			},
			"values": [
				[1504283307, "5.73"],
				[1504360587, "2.56"]
			]
		}, {
			"metric": {
				"host_name": "lhss-nyc-prod-4"
			},
			"values": [
				[1504293387, "7.2"],
				[1504313547, "2.33"],
				[1504323147, "2.81"],
				[1504335627, "3.54"],
				[1504358187, "3.44"]
			]
		}, {
			"metric": {
				"host_name": "lhss-nyc-prod-39"
			},
			"values": [
				[1504292907, "4.94"],
				[1504357707, "3.22"]
			]
		}, {
			"metric": {
				"host_name": "lhss-nyc-prod-40"
			},
			"values": [
				[1504280427, "5.26"],
				[1504293867, "4.89"]
			]
		}]
	}
} 
`)

	var o1 interface{}
	var o2 interface{}
	var err error
	z := new([]byte)

	err = MergeNaively(z, &a, &b)
	if err != nil {
		t.Error("mergeJson failed with", err)
	}

	json.Unmarshal(*z, &o1)
	json.Unmarshal(result, &o2)

	if !reflect.DeepEqual(o1, o2) {
		t.Error("Expected: ", string(result), "got", string(*z))
	}
}
