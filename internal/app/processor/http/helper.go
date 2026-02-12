package rprocessor

import (
	"net/http"

	"github.com/gorilla/mux"
)

//	func reg(r *mux.Router, method, path string, handler http.Handler) {
//		r.Methods(method).Path(path).Handler(handler)
//	}
func reg(r *mux.Router, path string, handler func(http.ResponseWriter, *http.Request), method string) {
	r.HandleFunc(path, handler).Methods(method)
}
