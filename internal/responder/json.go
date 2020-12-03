package responder

import (
	"encoding/json"
	"net/http"
)

//toJSON replies to the request with the specified payload and HTTP code
func toJSON(w http.ResponseWriter, httpCode int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")

	w.WriteHeader(httpCode)
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	if err := enc.Encode(payload); err != nil {
		http.Error(w, "Error encoding data", http.StatusInternalServerError)
	}
}
