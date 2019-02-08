package gopull

import (
	"log"
	"net/http"
	"time"
)

// Logger wraps an http.Handler and prints information about a request. This
// information includes a timestamp, the method, the endpoints, and the amount
// of time it took to complete the request.
func Logger(inner http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		inner.ServeHTTP(w, r)

		log.Printf(
			"%s\t%s\t%s\t%s",
			r.Method,
			r.RequestURI,
			name,
			time.Since(start),
		)
	})
}
