package api

import "net/http"

func newServeMux() *http.ServeMux {
	mux := http.NewServeMux()
	updateHandler := newUpdateHandler()
	mux.Handle("/update/", updateHandler)

	return mux
}
