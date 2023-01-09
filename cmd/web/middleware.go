package main

import (
	"net/http"

	"github.com/justinas/nosurf"
)

func SetupSession(next http.Handler)http.Handler{
	return sessionManager.LoadAndSave(next)
}
func NoSurf(next http.Handler) http.Handler {
	NoSurfHandler := nosurf.New(next)
	NoSurfHandler.SetBaseCookie(http.Cookie{
		Name:     "myCSRFCookie",
		Path:     "/",
		Domain:   "",
		Secure:   false,
		HttpOnly: true,
		MaxAge:   3600,
		SameSite: http.SameSiteLaxMode,
	})
	return NoSurfHandler
}