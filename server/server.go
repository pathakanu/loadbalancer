package server

import (
	"net/http"
)
// Server defines the interface for a backend server.
type Server interface {
	Address() string
	IsAlive() bool
	Serve(w http.ResponseWriter, r *http.Request)
}