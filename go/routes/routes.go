package routes

import (
	"encoding/json"
	"net/http"
)

type ResponseData struct {
	Message string `json:"message"`
}

func (r ResponseData) Write(w *http.ResponseWriter, code int) {
	data, _ := json.Marshal(r)
	(*w).Header().Set("Content-Type", "application/json")
	(*w).WriteHeader(code)
	(*w).Write(data)
}

func WriteJSON(code int, w *http.ResponseWriter, r any) {
	data, _ := json.Marshal(r)
	(*w).Header().Set("Content-Type", "application/json")
	(*w).WriteHeader(code)
	(*w).Write(data)
}
