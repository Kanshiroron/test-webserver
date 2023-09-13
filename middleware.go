package main

import (
	"net/http"

	log "github.com/sirupsen/logrus"
)

func LogMiddleware(downstream func(*log.Entry, http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		downstream(log.WithField(logEndpoint, r.URL.Path), w, r)
	}
}
