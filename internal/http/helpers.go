package httpserver

import (
	"AliceSkills/internal/http/middleware"
	"AliceSkills/internal/skills"
	"AliceSkills/internal/utils"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func ensurePost(w http.ResponseWriter, r *http.Request) bool {
	if r.Method == http.MethodPost {
		return true
	}
	w.Header().Set("Allow", http.MethodPost)
	w.WriteHeader(http.StatusMethodNotAllowed)
	return false
}

func ensureJSON(w http.ResponseWriter, r *http.Request) bool {
	ct := strings.ToLower(strings.TrimSpace(r.Header.Get("Content-Type")))
	if ct != "" && strings.HasPrefix(ct, "application/json") {
		return true
	}
	respondError(w, r, http.StatusUnsupportedMediaType, "unsupported_media_type", "Content-Type must be application/json")
	return false
}

func readBody(w http.ResponseWriter, r *http.Request, limit int64) ([]byte, bool) {
	r.Body = http.MaxBytesReader(w, r.Body, limit)

	body, err := io.ReadAll(r.Body)
	if err != nil {
		var maxErr *http.MaxBytesError
		if errors.As(err, &maxErr) {
			respondError(
				w, r,
				http.StatusRequestEntityTooLarge,
				"payload_too_large",
				fmt.Sprintf("payload too large (limit %d bytes)", maxErr.Limit),
			)
			return nil, false
		}
		respondError(w, r, http.StatusBadRequest, "bad_request", "failed to read request body")
		return nil, false
	}

	if len(bytes.TrimSpace(body)) == 0 {
		respondError(w, r, http.StatusBadRequest, "bad_request", "empty body")
		return nil, false
	}

	return body, true
}

func parseSkillFromPath(path string) (string, bool) {
	rest := strings.TrimPrefix(path, "/alice/")
	rest = strings.Trim(rest, "/")
	parts := strings.Split(rest, "/")
	if len(parts) == 2 && parts[0] != "" && parts[1] == "webhook" {
		return strings.ToLower(strings.TrimSpace(parts[0])), true
	}
	return "", false
}

func pingFastPath(w http.ResponseWriter, body []byte) bool {
	var input skills.Event
	_ = json.Unmarshal(body, &input)

	if strings.EqualFold(strings.TrimSpace(input.Request.OriginalUtterance), "ping") {
		utils.JsonHelper(w, http.StatusOK, &skills.Response{
			Version: input.Version,
			Session: input.Session,
			Result: struct {
				Text       string `json:"text"`
				EndSession bool   `json:"end_session"`
			}{
				Text:       "pong",
				EndSession: true,
			},
		})
		return true
	}
	return false
}

func respondError(w http.ResponseWriter, r *http.Request, status int, code, msg string) {
	rid := middleware.RequestID(r.Context())
	utils.JsonHelper(w, status, errorPayload(code, msg, rid))
}

func errorPayload(code, msg, rid string) map[string]any {
	return map[string]any{
		"error": map[string]any{
			"code":       code,
			"message":    msg,
			"request_id": rid,
		},
	}
}
