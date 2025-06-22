package middleware

import (
	"net/http"

	custErr "github.com/PIRSON21/mediasoft-go/internal/errors"
)

func Recoverer(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				// Следует использовать какой-то буферизованный ResponseWriter,
				// и если случилась паника, писать только это.
				// Может займусь этим.
				custErr.UnnamedError(w, http.StatusInternalServerError, "Internal Server Error")
			}
		}()

		next.ServeHTTP(w, r)
	}
}
