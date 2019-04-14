package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/google/jsonapi"
)

type Response struct {
	DataType string `jsonapi:"attr,data_type"`
}

func TestMyAPI(t *testing.T) {

	resultParser := func(payload string) string {
		r := strings.NewReader(payload)
		resp := &Response{}
		if err := jsonapi.UnmarshalPayload(r, resp); err != nil {
			t.Fatalf("invalid response payload, err: %v", err)
		}

		return resp.DataType
	}

	type payload struct {
		Event     string      `jsonapi:"primary,event"`
		EventData interface{} `jsonapi:"attr,data"`
	}

	tests := []struct {
		name                 string
		payload              payload
		wantHTTPResponseCode int
		wantResult           string
	}{
		{
			name: "test-1",
			payload: payload{
				Event: "event1",
				EventData: map[string]interface{}{
					"key": 1,
				},
			},
			wantHTTPResponseCode: 200,
			wantResult:           "int",
		},
		{
			name: "test-2",
			payload: payload{
				Event: "event2",
				EventData: map[string]interface{}{
					"key": 1.1,
				},
			},
			wantHTTPResponseCode: 200,
			wantResult:           "float64",
		},
		{
			name: "test-3",
			payload: payload{
				Event: "event3",
				EventData: map[string]interface{}{
					"key": "",
				},
			},
			wantHTTPResponseCode: 200,
			wantResult:           "string",
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			apiCall := func(m interface{}) (int, string) {

				pl, err := jsonapi.Marshal(m)
				if err != nil {
					t.Fatal(err)
				}

				plJSONBytes, err := json.Marshal(pl)
				if err != nil {
					t.Fatal(err)
				}

				resp, err := http.Post("http://localhost:8080", "", strings.NewReader(string(plJSONBytes)))
				if err != nil {
					t.Fatal(err)
				}
				defer resp.Body.Close()

				bodyBytes, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					t.Fatalf(err.Error())
				}

				return resp.StatusCode, resultParser(string(bodyBytes))

			}

			if httpResult, res := apiCall(&tt.payload); httpResult != tt.wantHTTPResponseCode || res != tt.wantResult {
				t.Errorf("httpResult = %v, want %v", httpResult, tt.wantHTTPResponseCode)
				t.Errorf("res = %v, want %v", res, tt.wantResult)
			}

		})

	}
}
