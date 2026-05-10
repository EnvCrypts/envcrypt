package server

import (
	"context"
	"log"
	"net"
	"net/http"

	"github.com/google/uuid"
	"github.com/vijayvenkatj/envcrypt/internal/errors"
	"github.com/vijayvenkatj/envcrypt/internal/helpers/reqcontext"
	"github.com/vijayvenkatj/envcrypt/internal/services"
)

type HandlerFunc func(w http.ResponseWriter, r *http.Request) error

func WithErrors(debug bool, next HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := next(w, r)
		if err == nil {
			return
		}

		reqID, ip, _ := reqcontext.GetRequestDetails(r.Context())
		log.Printf("request error: id=%s method=%s path=%s ip=%v err=%v", reqID, r.Method, r.URL.Path, ip, err)
		errors.Render(w, err, debug)
	}
}

func AuthMiddleware(sessionService *services.SessionService, next HandlerFunc) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		sessionID := r.Header.Get("X-Session-ID")
		if sessionID == "" {
			return errors.Unauthorized("SESSION_MISSING", "Session ID is required", "Log in to obtain a session")
		}

		sid, err := uuid.Parse(sessionID)
		if err != nil {
			return errors.Unauthorized("SESSION_INVALID", "Session ID format is invalid", "Provide a valid UUID session ID")
		}

		if err := sessionService.GetSession(r.Context(), sid); err != nil {
			return err
		}

		ctx := context.WithValue(r.Context(), "session_id", sid)
		return next(w, r.WithContext(ctx))
	}
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
