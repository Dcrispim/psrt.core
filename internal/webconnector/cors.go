package webconnector

import (
	"fmt"
	"net/http"
)

func WithStrictCORS(s *Server, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		isHealth := r.URL.Path == "/health"

		if origin != "" {
			if !s.IsOriginAllowed(origin) {
				s.audit.OriginRejected(origin, r.URL.Path)
				w.Header().Set("Content-Type", "text/plain; charset=utf-8")
				http.Error(
					w,
					fmt.Sprintf(
						"origin not allowed (got %s; configure allowed_origin in psrt-connector.ini, e.g. %s)",
						origin,
						origin,
					),
					http.StatusForbidden,
				)
				return
			}
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Vary", "Origin")
		} else if !isHealth {
			s.audit.OriginRejected("(missing)", r.URL.Path)
			http.Error(w, "origin required", http.StatusForbidden)
			return
		}

		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Header.Get("Access-Control-Request-Private-Network") == "true" {
			w.Header().Set("Access-Control-Allow-Private-Network", "true")
		}

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
