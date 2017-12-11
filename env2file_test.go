package main

import (
	"fmt"
	"reflect"
	"testing"
)

func Test_aggregateFromEnv(t *testing.T) {
	tests := []struct {
		variables []string
		wantTargs []target
		wantError []error
	}{
		{
			[]string{`EF_name_file=file.json`, `EF_data_file={"a":0}`},
			[]target{{name: "file.json", data: `{"a":0}`}},
			[]error{},
		},
		{
			[]string{`EF_name_file=fi_le.json`, `EF_data_file={"a":0}`},
			[]target{{name: "fi_le.json", data: `{"a":0}`}},
			[]error{},
		},
		{
			[]string{`EF_name_fi_le=fi_le.json`, `EF_data_fi_le={"a":0}`},
			[]target{{name: "fi_le.json", data: `{"a":0}`}},
			[]error{},
		},
	}
	for ii, tt := range tests {
		t.Run(fmt.Sprint(ii), func(t *testing.T) {
			gotTargs, gotError := aggregateFromEnv(tt.variables)
			equal(t, tt.wantTargs, gotTargs)
			equal(t, tt.wantError, gotError)
		})
	}
}

func Test_decodeKey(t *testing.T) {
	tests := []struct {
		key       string
		wantType  string
		wantName  string
		wantError error
	}{
		{"EF_name_target", "name", "target", nil},
		{"EF_data_target", "data", "target", nil},
		{"EF_name_target_underscore", "name", "target_underscore", nil},
		{"EF_name_", "name", "", nil},
		{"EF_name", "", "", fmt.Errorf("EF_name has invalid pattern")},
		{"EF_", "", "", fmt.Errorf("EF_ has invalid pattern")},
		{"EF", "", "", fmt.Errorf("EF has invalid pattern")},
	}
	for ii, tt := range tests {
		t.Run(fmt.Sprint(ii), func(t *testing.T) {
			gotType, gotName, gotError := decodeKey(tt.key)
			equal(t, tt.wantError, gotError)
			if gotError == nil {
				equal(t, tt.wantType, gotType)
				equal(t, tt.wantName, gotName)
			}
		})
	}
}

// helpers

func equal(t *testing.T, expected, actual interface{}) {
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("%v != %v", expected, actual)
	}
}
