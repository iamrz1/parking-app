package logger

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	ulid "github.com/oklog/ulid/v2"
	"github.com/rs/zerolog"
	"github.com/tylerb/gls"
)

// ReqIDTag holds tag for request id
const ReqIDTag = "RequestID"

type ridHook struct{}

// Run gets RequestID from middleware and sets it in the event context
func (h ridHook) Run(e *zerolog.Event, level zerolog.Level, msg string) {
	if id := gls.Get(ReqIDTag); id != nil && level != zerolog.NoLevel {
		e.Str(ReqIDTag, fmt.Sprintf("%v", id))
	}
}

// GenReqID generates a new request id for each request and stores in the gls. Storage is cleaned up after request is server.
func GenReqID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var entropy = ulid.Monotonic(rand.New(rand.NewSource(time.Now().UnixNano())), 0)
		gls.Set(ReqIDTag, ulid.MustNew(ulid.Timestamp(time.Now()), entropy))
		defer gls.Cleanup()
		next.ServeHTTP(w, r)
	})
}
