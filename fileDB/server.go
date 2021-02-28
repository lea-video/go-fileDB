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
	if r.Method == "GET" {
		resp = HandleGetRequest(ctx, w, r)
	} else {
		resp = &simpleFailed{
			err:        "Method is not supported",
			statusCode: http.StatusMethodNotAllowed,
		}
	}

	// send successful reply
	if !resp.DidFail() {
		bytes, err := json.Marshal(resp)
		if err != nil {
			// TODO: replace with logger
			fmt.Println(err)
			resp = &simpleFailed{
				err:        "Uups, sth went wrong",
				statusCode: http.StatusInternalServerError,
			}
		} else {
			w.WriteHeader(resp.GetStatusCode())
			// TODO: validate this
			// ignore errors
			// most likely connections problems
			// handled by connecting client
			_, _ = w.Write(bytes)
		}
	}
	// send failed reply
	// can get triggered from inside "send successful reply"
	if resp.DidFail() {
		http.Error(w, resp.GetError(), resp.GetStatusCode())
	}
}
