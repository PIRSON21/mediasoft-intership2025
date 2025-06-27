package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/PIRSON21/mediasoft-go/pkg/logger"
	"go.uber.org/zap"
)

func LoggingMiddleware(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := logger.GetLogger().With(
			zap.String("request-id", GetRequestID(r.Context())),
		)
		start := time.Now()

		next.ServeHTTP(w, r)

		// Хотелось бы, как в норм роутерах, вывод ещё статуса ответа,
		// но в net/http нельзя обратиться к ResponseWriter, чтобы получить статус ответа.
		// Можно было сделать обертку ResponseWriter, где ещё хранить статус, а здесь его выводить.
		// Может займусь.
		log.Info(fmt.Sprintf("%s %s", r.Method, time.Since(start)))
	}
}
