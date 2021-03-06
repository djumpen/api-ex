package main

import (
	"log"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/google/jsonapi"
)

type TypeReq struct {
	Event     string                 `jsonapi:"primary,event"`
	EventData map[string]interface{} `jsonapi:"attr,data"`
}

type TypeResp struct {
	DataType string `jsonapi:"attr,data_type"`
}

func main() {
	http.HandleFunc("/", typeHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// typeHandler is a handler which returns type of incoming value
func typeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", jsonapi.MediaType)

	pl := new(TypeReq)

	if err := jsonapi.UnmarshalPayload(r.Body, pl); err != nil {
		jsonapi.MarshalErrors(w, []*jsonapi.ErrorObject{{
			Title:  "Validation Error",
			Detail: err.Error(),
			Status: "400",
		}})
		return
	}

	k, ok := pl.EventData["key"]
	if !ok {
		jsonapi.MarshalErrors(w, []*jsonapi.ErrorObject{{
			Title:  "Validation Error",
			Detail: "Given request body was invalid.",
			Status: "400",
		}})
		return
	}

	dt := reflect.TypeOf(k).String()

	if kf, ok := k.(float64); ok {
		if !hasDecimals(kf) {
			dt = "int"
		}
	}

	resp := &TypeResp{
		DataType: dt,
	}

	w.WriteHeader(http.StatusOK)

	if err := jsonapi.MarshalPayload(w, resp); err != nil {
		jsonapi.MarshalErrors(w, []*jsonapi.ErrorObject{{
			Title:  "Internal Error",
			Detail: err.Error(),
			Status: "500",
		}})
	}
}

// hasDecimals shows the fact of existence of decimals in a given float64
// If they are presented - true will be returned
func hasDecimals(v float64) bool {
	s := strconv.FormatFloat(v, 'f', -1, 64)
	i := strings.IndexByte(s, '.')
	return i > -1
}
