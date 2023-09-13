package main

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
)

const (
	serverQueryParamCrashExitCode string = "exit_code"
)

// echo returns the request body
func crash(l *log.Entry, w http.ResponseWriter, r *http.Request) {
	// default exit code
	exitCode := 1

	// parsing query
	exitCodeString := r.URL.Query().Get(serverQueryParamCrashExitCode)
	if len(exitCodeString) > 0 {
		var err error
		if exitCode, err = strconv.Atoi(exitCodeString); err != nil {
			l.WithError(err).Errorf("failed to parse %q query param", serverQueryParamCrashExitCode)
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "failed to parse %q query param: %s", serverQueryParamCrashExitCode, err.Error())
			return
		}
	}
	w.WriteHeader(http.StatusOK)
	l.Infof("server is about to crash with exit code: %d", exitCode)

	// crash
	go func(exitCode int) {
		time.Sleep(time.Second)
		syscall.Exit(exitCode)
	}(exitCode)
}

// echo returns the request body
func echo(l *log.Entry, w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	io.Copy(w, r.Body)
}
