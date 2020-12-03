package responder

import (
	"encoding/xml"
	"net/http"
)

//toXML replies to the request with the specified payload and HTTP code
func toXML(w http.ResponseWriter, httpCode int, payload interface{}) {
	w.Header().Set("Content-Type", "application/xml; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")

	w.WriteHeader(httpCode)
	enc := xml.NewEncoder(w)
	if err := enc.Encode(payload); err != nil {
		http.Error(w, "Error encoding data", http.StatusInternalServerError)
	}
}
