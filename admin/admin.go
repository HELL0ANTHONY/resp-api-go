package admin

import (
	"net/http"
	"os"
)

type adminPortal struct {
	password string
}

func NewAdminPortal() *adminPortal {
	password := os.Getenv("ADMIN_PASSWORD")
	if password == "" {
		panic("required env var ADMIN_PASSWORD not set")
	}
	return &adminPortal{password: password}
}

func (a adminPortal) Handler(w http.ResponseWriter, r *http.Request) {
	user, pass, ok := r.BasicAuth()
	if !ok || user != "admin" || pass != a.password {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte("401 - unauthorized"))
		return
	}
	_, _ = w.Write([]byte("<html><h1>Super secret admin portal</h1></html>"))
}
