package middleware

import (
	"AliceSkills/internal/utils"
	"net/http"
)

const (
	InternalErrorCode    = "internal_error"
	InternalErrorMessage = "internal server error"
)

type ErrorResponse struct {
	RequestID string `json:"-"`
}

func (e ErrorResponse) JSON() map[string]any {
	return map[string]any{
		"error": map[string]string{
			"code":       InternalErrorCode,
			"message":    InternalErrorMessage,
			"request_id": e.RequestID,
		},
	}
}

func RecoverMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				rid := RequestID(r.Context())
				if rid == "" {
					rid = w.Header().Get(headerRequestID)
				}
				payload := ErrorResponse{RequestID: rid}.JSON()
				w.Header().Set("X-Request-ID", rid)
				utils.JsonHelper(w, http.StatusInternalServerError, payload)
				return
			}
		}()
		next.ServeHTTP(w, r)
	})
}
