package middleware

import (
	"log"
	"net/http"
	"time"
)

// responseWriter обёртка для записи статус кода
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// Logger middleware для логирования HTTP запросов
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Создаём обёртку для response writer
		wrapped := &responseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		// Выполняем запрос
		next.ServeHTTP(wrapped, r)

		// Логируем информацию о запросе
		duration := time.Since(start)
		log.Printf(
			"[%s] %s %s - Status: %d - Duration: %v",
			r.Method,
			r.RequestURI,
			r.RemoteAddr,
			wrapped.statusCode,
			duration,
		)
	})
}
