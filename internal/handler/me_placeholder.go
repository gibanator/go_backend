package handler

import "net/http"

func WriteMePlaceholder(w http.ResponseWriter, userID string) {
	writeJSON(w, http.StatusOK, map[string]string{
		"id": userID,
	})
}
