package secrets

import (
		"net/http"
)

func headers(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Add("Content-Type", "application/json")
				w.Header().Add("Access-Control-Allow-Origin", "*")
				w.Header().Add("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
				w.Header().Add("Access-Control-Allow-Headers", "Origin, Content-Type")

				if r.Method == "OPTIONS" {
						return
				}

				h.ServeHTTP(w, r)
		})
}


func NewHttpServer(endpoints Endpoints) http.Handler {
		s := http.NewServeMux()
		r := MakeHandler(endpoints)

		r.Use(headers)

		s.Handle("/", r)
		return s
}
