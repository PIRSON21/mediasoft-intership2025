package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

type RequestIDType string

var requestIDKey RequestIDType = "x-request-id"

// RequestID добавляет уникальный идентификатор запроса в контекст запроса.
func RequestID(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		requestID := r.Header.Get("x-request-id")
		if requestID == "" {
			requestID = createRequestID()
		}

		r = r.WithContext(context.WithValue(r.Context(), requestIDKey, requestID))

		next.ServeHTTP(w, r)
	}
}

// createRequestID генерирует новый уникальный идентификатор запроса.
func createRequestID() string {
	return uuid.NewString()
}

// GetRequestID извлекает уникальный идентификатор запроса из контекста.
func GetRequestID(requestCtx context.Context) string {
	v := requestCtx.Value(requestIDKey)
	if v == nil {
		return ""
	}
	return v.(string)
}
