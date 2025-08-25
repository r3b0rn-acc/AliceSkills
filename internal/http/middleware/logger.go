package middleware

import (
	"log"
	"net/http"
	"time"
)

type responseRecorder struct {
	http.ResponseWriter
	status      int
	bytes       int
	wroteHeader bool
}

func (rr *responseRecorder) WriteHeader(code int) {
	if rr.wroteHeader {
		return
	}
	rr.status = code
	rr.wroteHeader = true
	rr.ResponseWriter.WriteHeader(code)
}

func (rr *responseRecorder) Write(b []byte) (int, error) {
	if !rr.wroteHeader {
		rr.WriteHeader(http.StatusOK)
	}
	n, err := rr.ResponseWriter.Write(b)
	rr.bytes += n
	return n, err
}

func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rr := &responseRecorder{ResponseWriter: w}

		next.ServeHTTP(rr, r)

		dur := time.Since(start)
		rid := RequestID(r.Context())
		if rid == "" {
			rid = w.Header().Get(headerRequestID)
		}
		method := r.Method
		path := r.URL.RequestURI()
		status := rr.status
		if status == 0 {
			status = http.StatusOK
		}
		bytes := rr.bytes

		ua := r.UserAgent()
		remote := r.Header.Get("X-Forwarded-For")
		if remote == "" {
			remote = r.RemoteAddr
		}

		log.Printf(
			"method=%s path=%s status=%d dur=%s bytes=%d rid=%s ua=%q remote=%s",
			method, path, status, dur, bytes, rid, ua, remote,
		)
	})
}
