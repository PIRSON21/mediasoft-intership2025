package middleware

import (
	"net/http"

	custErr "github.com/PIRSON21/mediasoft-intership2025/internal/errors"
	"github.com/PIRSON21/mediasoft-intership2025/pkg/logger"
	"go.uber.org/zap"
)

func Recoverer(next http.Handler) http.HandlerFunc {
	log := logger.GetLogger()
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				// Следует использовать какой-то буферизованный ResponseWriter,
				// и если случилась паника, писать только это.
				// Может займусь этим.
				log.Error("panic", zap.Any("err", r))
				custErr.UnnamedError(w, http.StatusInternalServerError, "Internal Server Error")
			}
		}()

		next.ServeHTTP(w, r)
	}
}
