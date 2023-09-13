package main

import (
	"strconv"
	"syscall"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

const (
	// environment
	envDebug              string = "DEBUG"
	envHTTPServerPort     string = "HTTP_PORT"
	envHTTPMonitoringPort string = "HTTP_MONITORING_PORT"

	// defaults
	defaultHTTPServerPort     int = 8080
	defaultHTTPMonitoringPort int = 8081
)

type Config struct {
	Debug              bool
	HTTPServerPort     int
	HTTPMonitoringPort int

	MonitoringEndpointsConfig MonitoringEndpointsConfig
}

func ParseConfigFromEnv() (c Config, err error) {
	// debug
	if debugLogString, found := syscall.Getenv(envDebug); found {
		if debugLog, err := strconv.ParseBool(debugLogString); err != nil {
			return Config{}, errors.WithMessagef(err, "failed to parse the debug value to boolean (env variable: %s)", envDebug)
		} else if debugLog {
			log.SetLevel(log.DebugLevel)
			log.Debug("debug logs enabled")
		}
	}

	// http server port
	if serverPortString, found := syscall.Getenv(envHTTPServerPort); found {
		if serverPort, err := strconv.Atoi(serverPortString); err != nil {
			return Config{}, errors.WithMessagef(err, "failed to parse server port to int (env variable: %s)", envHTTPServerPort)
		} else {
			c.HTTPServerPort = serverPort
		}
	} else {
		log.Debugf("using default server port: %d", defaultHTTPServerPort)
		c.HTTPServerPort = defaultHTTPServerPort
	}

	// http monitoring port
	if monitoringPortString, found := syscall.Getenv(envHTTPMonitoringPort); found {
		if monitoringPort, err := strconv.Atoi(monitoringPortString); err != nil {
			return Config{}, errors.WithMessagef(err, "failed to parse monitoring port to int (env variable: %s)", envHTTPMonitoringPort)
		} else {
			c.HTTPMonitoringPort = monitoringPort
		}
	} else {
		log.Debugf("using default monitoring port: %d", defaultHTTPMonitoringPort)
		c.HTTPMonitoringPort = defaultHTTPMonitoringPort
	}

	// monitoring endpoints config
	c.MonitoringEndpointsConfig, err = ParseMonitoringEndpointsConfigFromEnv()
	if err != nil {
		return Config{}, err
	}
	return c, nil
}
