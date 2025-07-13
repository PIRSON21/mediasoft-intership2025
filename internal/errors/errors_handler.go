package errors

import (
	"encoding/json"
	"errors"
	"net/http"
)

// UnnamedError отправляет ошибку в формате JSON с заданным статус-кодом
// в формате { "error": "message" } или { "error": ["message1", "message2"] }.
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

// Any проверяет, содержится ли ошибка в списке целевых ошибок.
// Возвращает true, если ошибка найдена среди целевых ошибок, иначе false.
func Any(err error, targets ...error) bool {
	for _, target := range targets {
		if errors.Is(err, target) {
			return true
		}
	}

	return false
}
