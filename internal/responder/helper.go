package responder

import "net/http"

const xmlContent = "application/xml"

//Response is responsible for calling proper method based on accepted response content
func Response(w http.ResponseWriter, r *http.Request, code int, payload interface{}) {
	acceptedContent := r.Header.Get("Accept")
	if acceptedContent == xmlContent {
		toXML(w, code, payload)
	} else {
		toJSON(w, code, payload)
	}
}

//SecretError is for service errors
type SecretError struct {
	Description string `json:"description" xml:"description"`
}

//WithError responds with secretError message
func WithError(w http.ResponseWriter, r *http.Request, code int, message string) {
	msg := SecretError{Description: message}
	Response(w, r, code, msg)
}
