package health

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"context"
)

// Handler returns an http.HandlerFunc that reports service health.
// It pings the database and returns a JSON status response.
func Handler(db *sql.DB, service, version string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()

		if err := db.PingContext(ctx); err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			json.NewEncoder(w).Encode(map[string]string{
				"status":  "error",
				"service": service,
				"version": version,
			})
			return
		}

		json.NewEncoder(w).Encode(map[string]string{
			"status":  "ok",
			"service": service,
			"version": version,
		})
	}
}
