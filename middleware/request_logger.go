package middleware

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"supertal-tha-parking-app/logger"
)

// RequestLogger ...
func RequestLogger() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger.GetLogger().Infof("Request: [%s] %s%s", r.Method, r.Host, r.URL)
			if body, err := readIntact(r); err == nil {
				logger.GetLogger().Debugf("\n%v\n", string(body))
			} else {
				logger.GetLogger().Error(err.Error())
			}

			next.ServeHTTP(w, r)
		})
	}
}

func readIntact(r *http.Request) ([]byte, error) {
	var buf bytes.Buffer
	tee := io.TeeReader(r.Body, &buf)
	body, err := ioutil.ReadAll(tee)
	r.Body = ioutil.NopCloser(&buf)
	return body, err
}
