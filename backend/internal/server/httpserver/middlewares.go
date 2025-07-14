package httpserver

import "net/http"

func CORSMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*") // or "http://localhost:3000"
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

        // Handle preflight request
        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusOK)
            return
        }

		h.ServeHTTP(w, r)
	})
}