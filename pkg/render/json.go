package render

import (
	"encoding/json"
	"net/http"
)

func JSON(w http.ResponseWriter, statusCode int, response any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	return json.NewEncoder(w).Encode(response)
}
