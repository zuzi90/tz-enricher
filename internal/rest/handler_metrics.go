package rest

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

func (s *Server) metrics(w http.ResponseWriter, r *http.Request) {
	prom := promhttp.Handler()
	prom.ServeHTTP(w, r)
	return
}
