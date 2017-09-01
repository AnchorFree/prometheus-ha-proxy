package merger

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestSimpleJsonMerge(t *testing.T) {
	a := []byte(`{"status":"success","data":{"resultType":"vector","result":[{"metric":{"host_name":"hss4"},"value":[1504157313.787,"24"]},{"metric":{"host_name":"hss508"},"value":[1504157313.787,"32"]}]}}`)
	b := []byte(`{"status":"success","data":{"resultType":"vector","result":[{"metric":{"host_name":"hss4"},"value":[1504157435.713,"20"]},{"metric":{"host_name":"bhss-nyc2-prod-80"},"value":[1504157435.713,"2"]}]}}`)
	//c := []byte(`{"status":"success","data":{"resultType":"vector","result":[{"metric":{"host_name":"hss4"},"value":[1504157435.713,"20"]},{"metric":{"host_name":"bhss-nyc2-prod-80"},"value":[1504157435.713,"2"]}]}}`)
	result := []byte(`{"status":"success","data":{"resultType":"vector","result":[{"metric":{"host_name":"hss4"},"value":[1504157313.787,"24"]},{"metric":{"host_name":"hss508"},"value":[1504157313.787,"32"]},{"metric":{"host_name":"hss4"},"value":[1504157435.713,"20"]},{"metric":{"host_name":"bhss-nyc2-prod-80"},"value":[1504157435.713,"2"]}]}}`)

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

func TestSimpleJsonMergeVariableArgs(t *testing.T) {
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
