package handlers

import (
	"encoding/json"
	"net/http"
	"time"
)

type HelloResponse struct {
	Message string `json:"message"`
}

func HelloHandler(w http.ResponseWriter, r *http.Request) {
	time.Sleep(2 * time.Second)
	response := HelloResponse{
		Message: "Hello",
	}
	json.NewEncoder(w).Encode(response)
	return
}
