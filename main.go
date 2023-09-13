package main

import (
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
)

func main() {
	// parsing config
	config, err := ParseConfigFromEnv()
	if err != nil {
		log.WithError(err).Fatal("failed to parse configuration from environment")
	}

	// server endpoints
	serverMux := http.NewServeMux()
	serverMux.HandleFunc("/crash", LogMiddleware(crash))
	serverMux.HandleFunc("/echo", LogMiddleware(echo))

	// monitoring endpoints
	monitoringEndpoints := NewMonitoringEndpoints(config.MonitoringEndpointsConfig)
	monitoringMux := http.NewServeMux()
	monitoringMux.HandleFunc("/liveness", LogMiddleware(monitoringEndpoints.Liveness))
	monitoringMux.HandleFunc("/readiness", LogMiddleware(monitoringEndpoints.Readiness))

	// starting webserver
	webserverListen := ":" + strconv.Itoa(config.HTTPServerPort)
	webserver := &http.Server{
		Handler: serverMux,
		Addr:    webserverListen,
	}
	go func() {
		err = webserver.ListenAndServe()
	}()
	time.Sleep(time.Second)
	if err != nil {
		log.WithError(err).Fatal("failed to start webserver")
	} else {
		log.Infof("webserver listening on %s", webserverListen)
	}

	// starting monitoring
	monitoringListen := ":" + strconv.Itoa(config.HTTPMonitoringPort)
	monitoring := &http.Server{
		Handler: monitoringMux,
		Addr:    monitoringListen,
	}
	go func() {
		err = monitoring.ListenAndServe()
	}()
	time.Sleep(time.Second)
	if err != nil {
		log.WithError(err).Fatal("failed to start monitoring webserver")
	} else {
		log.Infof("monitoring webserver listening on %s", monitoringListen)
	}

	log.Info("test webserver is up and running")

	// waiting for stop signal
	stopSignal := make(chan os.Signal, 1)
	signal.Notify(stopSignal, os.Interrupt)
	<-stopSignal
	log.Info("stop signal received, stopping server")
	webserver.Close()
	log.Info("test webserver stopped")
}
