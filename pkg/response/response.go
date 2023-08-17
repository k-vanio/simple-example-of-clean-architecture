package response

import (
	"encoding/json"
	"net/http"
)

const (
	// types
	ContentTypeJson = "application/json"
)

func Json(w http.ResponseWriter, status int, body interface{}) {
	w.Header().Set("Content-Type", ContentTypeJson)
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(body); err != nil {
		Json(w, http.StatusInternalServerError, struct {
			Message string `json:"response"`
		}{Message: "internal server error"})
	}
}
