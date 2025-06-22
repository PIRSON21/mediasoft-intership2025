package errors

import (
	"encoding/json"
	"net/http"
)

func UnnamedError(w http.ResponseWriter, statusCode int, errMessages ...string) error {
	var response any
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if len(errMessages) == 1 {
		response = map[string]string{
			"error": errMessages[0],
		}
	} else {
		response = map[string][]string{
			"error": errMessages,
		}
	}

	return json.NewEncoder(w).Encode(response)
}
