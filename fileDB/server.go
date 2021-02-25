package fileDB

import "net/http"

func RegisterServer() error {
	http.HandleFunc("/", reqHandler)
	return http.ListenAndServe(":8080", nil)
}

func reqHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		HandleGetRequest(w, r)
	} else {
		denie(w, r)
	}
}

func denie(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Method is not supported", http.StatusMethodNotAllowed)
}