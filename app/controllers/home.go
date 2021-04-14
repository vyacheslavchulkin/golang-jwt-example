package controllers

import (
	"fmt"
	"net/http"
)

func Home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		notFound(w)
		return
	}
	w.WriteHeader(200)
	fmt.Fprintf(w, "Welcome!")
}
