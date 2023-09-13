package main

import (
	"net/http"
	"strconv"
	"syscall"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

const (
	// environment
	envMonitoringLivenessStatusOk     string = "MONITORING_LIVENESS_STATUS_OK"
	envMonitoringLivenessStatusError  string = "MONITORING_LIVENESS_STATUS_ERROR"
	envMonitoringReadinessStatusOk    string = "MONITORING_READINESS_STATUS_OK"
	envMonitoringReadinessStatusError string = "MONITORING_READINESS_STATUS_ERROR"

	// defaults
	defaultMonitoringLivenessStatusOk     int = http.StatusOK
	defaultMonitoringLivenessStatusError  int = http.StatusInternalServerError
	defaultMonitoringReadinessStatusOk    int = http.StatusOK
	defaultMonitoringReadinessStatusError int = http.StatusInternalServerError
)

type MonitoringEndpointsConfig struct {
	LivenessStatusOk     int
	LivenessStatusError  int
	ReadinessStatusOk    int
	ReadinessStatusError int
}

func ParseMonitoringEndpointsConfigFromEnv() (c MonitoringEndpointsConfig, err error) {
	// liveness status ok
	if livenessStatusOkString, found := syscall.Getenv(envMonitoringLivenessStatusOk); found {
		if livenessStatusOk, err := strconv.Atoi(livenessStatusOkString); err != nil {
			return MonitoringEndpointsConfig{}, errors.WithMessagef(err, "failed to parse liveness status ok to int (env variable: %s)", envMonitoringLivenessStatusOk)
		} else {
			c.LivenessStatusOk = livenessStatusOk
		}
	} else {
		log.Debugf("using default liveness status ok: %d", defaultMonitoringLivenessStatusOk)
		c.LivenessStatusOk = defaultMonitoringLivenessStatusOk
	}
	// liveness status error
	if livenessStatusErrorString, found := syscall.Getenv(envMonitoringLivenessStatusError); found {
		if livenessStatusError, err := strconv.Atoi(livenessStatusErrorString); err != nil {
			return MonitoringEndpointsConfig{}, errors.WithMessagef(err, "failed to parse liveness status ok to int (env variable: %s)", envMonitoringLivenessStatusError)
		} else {
			c.LivenessStatusError = livenessStatusError
		}
	} else {
		log.Debugf("using default liveness status ok: %d", defaultMonitoringLivenessStatusError)
		c.LivenessStatusError = defaultMonitoringLivenessStatusError
	}
	// readiness status ok
	if readinessStatusOkString, found := syscall.Getenv(envMonitoringReadinessStatusOk); found {
		if readinessStatusOk, err := strconv.Atoi(readinessStatusOkString); err != nil {
			return MonitoringEndpointsConfig{}, errors.WithMessagef(err, "failed to parse liveness status ok to int (env variable: %s)", envMonitoringReadinessStatusOk)
		} else {
			c.ReadinessStatusOk = readinessStatusOk
		}
	} else {
		log.Debugf("using default liveness status ok: %d", defaultMonitoringLivenessStatusOk)
		c.ReadinessStatusOk = defaultMonitoringLivenessStatusOk
	}
	// readiness status error
	if readinessStatusErrorString, found := syscall.Getenv(envMonitoringReadinessStatusError); found {
		if readinessStatusError, err := strconv.Atoi(readinessStatusErrorString); err != nil {
			return MonitoringEndpointsConfig{}, errors.WithMessagef(err, "failed to parse liveness status ok to int (env variable: %s)", envMonitoringReadinessStatusError)
		} else {
			c.ReadinessStatusError = readinessStatusError
		}
	} else {
		log.Debugf("using default liveness status ok: %d", defaultMonitoringReadinessStatusError)
		c.ReadinessStatusError = defaultMonitoringReadinessStatusError
	}
	return c, nil
}
