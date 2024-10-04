package server

import (
	"log/slog"
	"net"
	"net/http"
)

func (srv *Server) logRequests() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userAttrs := slog.Group("user", "ip", r.RemoteAddr)
			addr, _, _ := net.SplitHostPort(r.RemoteAddr)
			requestAttrs := slog.Group("request",
				"content_length", r.ContentLength,
				"method", r.Method,
				"proto", r.Proto,
				"remote_addr", addr,
				"uri", r.RequestURI,
				"user_agent", r.UserAgent(),
			)
			srv.logger.Info("received request", userAttrs, requestAttrs)
			next.ServeHTTP(w, r)
		})
	}
}
