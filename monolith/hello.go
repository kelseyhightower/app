package main

import (
	"net/http"
	"time"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	time.Sleep(2 * time.Second)
	w.Write([]byte("<h1>Hello</h1>"))
	return
}
