package server

import (
	"context"
	"net"
	"net/http"

	"github.com/google/uuid"
	"github.com/vijayvenkatj/envcrypt/internal/helpers"
	"github.com/vijayvenkatj/envcrypt/internal/helpers/reqcontext"
	"github.com/vijayvenkatj/envcrypt/internal/services"
)

func AuthMiddleware(sessionService *services.SessionService, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		sessionID := r.Header.Get("X-Session-ID")
		if sessionID == "" {
			helpers.WriteError(w, http.StatusUnauthorized, "Session ID is required")
			return
		}

		sid, err := uuid.Parse(sessionID)
		if err != nil {
			helpers.WriteError(w, http.StatusUnauthorized, "Session ID is invalid")
			return
		}

		if err := sessionService.GetSession(r.Context(), sid); err != nil {
			helpers.WriteError(w, http.StatusUnauthorized, "Session ID is invalid or expired")
			return
		}

		ctx := context.WithValue(r.Context(), "session_id", sid)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func RequestMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := r.Header.Get("X-Request-Id")
		if reqID == "" {
			reqID = uuid.New().String()
		}

		ip := r.Header.Get("X-Forwarded-For")
		if ip == "" {
			host, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				ip = r.RemoteAddr
			} else {
				ip = host
			}
		}

		ua := r.UserAgent()

		ctx := reqcontext.SetRequestDetails(r.Context(), reqID, ip, ua)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
