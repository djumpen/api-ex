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

func typeHandler(w http.ResponseWriter, r *http.Request) {
	pl := new(TypeReq)

	if err := jsonapi.UnmarshalPayload(r.Body, pl); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	k, ok := pl.EventData["key"]
	if !ok {
		w.Header().Set("Content-Type", jsonapi.MediaType)
		jsonapi.MarshalErrors(w, []*jsonapi.ErrorObject{{
			Title:  "Validation Error",
			Detail: "Given request body was invalid.",
			Status: "400",
		}})
		return
	}

	var dt string

	switch k := k.(type) {
	case float64:
		dt = "int"
		if hasDecimals(k) {
			dt = "float64"
		}
	default:
		dt = reflect.TypeOf(k).String()
	}

	resp := &TypeResp{
		DataType: dt,
	}

	w.Header().Set("Content-Type", jsonapi.MediaType)
	w.WriteHeader(http.StatusOK)

	if err := jsonapi.MarshalPayload(w, resp); err != nil {
		jsonapi.MarshalErrors(w, []*jsonapi.ErrorObject{{
			Title:  "Internal Error",
			Detail: err.Error(),
			Status: "500",
		}})
	}
}

func hasDecimals(v float64) bool {
	s := strconv.FormatFloat(v, 'f', -1, 64)
	i := strings.IndexByte(s, '.')
	return i > -1
}
