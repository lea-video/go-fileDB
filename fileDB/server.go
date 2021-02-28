package fileDB

import (
	"net/http"
)

func RegisterServer(ctx Context) error {
	calReqHandler := func(w http.ResponseWriter, r *http.Request) { reqHandler(ctx, w, r) }
	http.HandleFunc("/", calReqHandler)
	return http.ListenAndServe(ctx.GetPortStr(), nil)
}

func reqHandler(ctx Context, w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		HandleGetRequest(ctx, w, r)
	} else {
		denie(w, r)
	}
}

func denie(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Method is not supported", http.StatusMethodNotAllowed)
}