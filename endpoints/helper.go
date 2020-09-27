package endpoints

import (
	"encoding/json"
	"encoding/xml"
	"net/http"
)

func Response(w http.ResponseWriter, r *http.Request, message interface{}, code int) {
	switch a := r.Header.Get("Accept"); {
	case a == "application/json":
		errorHandler(toJSON, w, message, code)
	case a == "application/xml":
		errorHandler(toXML, w, message, code)
	}
}

func toJSON(w http.ResponseWriter, content interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(content)
}

func toXML(w http.ResponseWriter, content interface{}) error {
	w.Header().Set("Content-Type", "application/xml")
	return xml.NewEncoder(w).Encode(content)
}

type ParseFunc func(w http.ResponseWriter, content interface{}) error

func errorHandler(f ParseFunc, w http.ResponseWriter, message interface{}, code int) {
	if err := f(w, message); err != nil {
		http.Error(w, err.Error(), code)
	}
}
