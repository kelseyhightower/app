package handlers

import (
	"net/http"
	"time"
)

func HelloHandler(w http.ResponseWriter, r *http.Request) {
	time.Sleep(2 * time.Second)
	w.Write([]byte("<h1>Hello</h1>"))
	return
}
