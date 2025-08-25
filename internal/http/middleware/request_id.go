package middleware

import (
	"context"
	"crypto/rand"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type ctxKey string

const ctxKeyRequestID ctxKey = "request_id"
const headerRequestID = "X-Request-ID"

func RequestID(ctx context.Context) string {
	val, _ := ctx.Value(ctxKeyRequestID).(string)
	return val
}

func RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rid := strings.TrimSpace(r.Header.Get(headerRequestID))
		if rid == "" {
			if v4, err := NewUUIDv4(); err == nil {
				rid = v4
			} else {
				rid = strconv.FormatInt(time.Now().UnixNano(), 16)
			}
		}

		w.Header().Set(headerRequestID, rid)

		ctx := context.WithValue(r.Context(), ctxKeyRequestID, rid)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func NewUUIDv4() (string, error) {
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		return "", err
	}

	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80

	s := fmt.Sprintf("%02x%02x%02x%02x-%02x%02x-%02x%02x-%02x%02x-%02x%02x%02x%02x%02x%02x",
		b[0], b[1], b[2], b[3],
		b[4], b[5],
		b[6], b[7],
		b[8], b[9],
		b[10], b[11], b[12], b[13], b[14], b[15],
	)

	return s, nil
}
