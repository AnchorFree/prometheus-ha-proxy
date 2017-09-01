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
				[1504286224, "0.12527084236321268"]
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
				[1504286220, "0.12527084236321268"],
				[1504286224, "0.12527084236321268"]
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
