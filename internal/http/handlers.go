package httpserver

import (
	"AliceSkills/internal/skills"
	"AliceSkills/internal/utils"
	"net/http"
)

const maxBodyBytes = 1 << 20

func aliceWebhookHandler(reg skills.Registry) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !ensurePost(w, r) {
			return
		}
		if !ensureJSON(w, r) {
			return
		}

		skillName, ok := parseSkillFromPath(r.URL.Path)
		if !ok {
			respondError(w, r, http.StatusNotFound, "not_found", "route not found")
			return
		}

		skill, ok := reg.Get(skillName)
		if !ok {
			respondError(w, r, http.StatusNotFound, "skill_not_found", "unknown skill")
			return
		}

		body, ok := readBody(w, r, maxBodyBytes)
		if !ok {
			return
		}

		if pingFastPath(w, body) {
			return
		}

		resp, err := skill.Handle(r.Context(), body)
		if err != nil {
			respondError(w, r, http.StatusInternalServerError, "internal_error", "internal server error")
			return
		}

		utils.JsonHelper(w, http.StatusOK, resp)
	}
}
