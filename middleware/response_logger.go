package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"supertal-tha-parking-app/logger"
)

// ResponseLogger ...
func ResponseLogger() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			res := httptest.NewRecorder()
			next.ServeHTTP(res, r)
			logger.GetLogger().Debugf("Response: %v", strings.TrimSpace(res.Body.String()))

			for k, v := range res.Header() {
				w.Header()[k] = v
			}
			w.WriteHeader(res.Code)
			res.Body.WriteTo(w)
		})
	}
}
