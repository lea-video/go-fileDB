package fileDB

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type serverResponse interface {
	GetStatusCode() int
	DidFail() bool
	GetError() string
}

type simpleFailed struct {
	err string
	statusCode int
}
func (simpleFailed) DidFail() bool { return true }
func (sf *simpleFailed) GetStatusCode() int { return sf.statusCode }
func (sf *simpleFailed) GetError() string { return sf.err }

func RegisterServer(ctx Context) error {
	calReqHandler := func(w http.ResponseWriter, r *http.Request) { reqHandler(ctx, w, r) }
	http.HandleFunc("/", calReqHandler)
	return http.ListenAndServe(ctx.GetPortStr(), nil)
}

func reqHandler(ctx Context, w http.ResponseWriter, r *http.Request) {
	var resp serverResponse

	// build reply
	if r.Method == http.MethodGet {
		resp = onGetRequest(ctx, r)
	} else if r.Method == http.MethodDelete {
		resp = onDeleteRequest(ctx, r)
	} else if r.Method == http.MethodPost {
		resp = onPostRequest(ctx, r)
	} else {
		resp = &simpleFailed{"Method is not supported", http.StatusMethodNotAllowed}

	}

	// send successful reply
	if !resp.DidFail() {
		bytes, err := json.Marshal(resp)
		if err != nil {
			// TODO: replace with logger
			fmt.Println(err)
			resp = &simpleFailed{"Uups, sth went wrong", http.StatusInternalServerError}
		} else {
			w.WriteHeader(resp.GetStatusCode())
			// TODO: validate this
			// ignore errors here
			// most likely connections problems
			// handled by connecting client
			_, err = w.Write(bytes)
			fmt.Println(err)
		}
	}
	// send failed reply
	// can get triggered from inside "send successful reply"
	if resp.DidFail() {
		http.Error(w, resp.GetError(), resp.GetStatusCode())
	}
}
