package main

import (
	"fmt"
	"net/http"
	"strconv"
	"sync"

	log "github.com/sirupsen/logrus"
)

const (
	monitoringQueryParamFail string = "fail"
)

func NewMonitoringEndpoints(config MonitoringEndpointsConfig) *MonitoringEndpoints {
	return &MonitoringEndpoints{
		lock:   &sync.Mutex{},
		config: config,
	}
}

type MonitoringEndpoints struct {
	lock            *sync.Mutex
	config          MonitoringEndpointsConfig
	failReadiness   bool
	failReadinessNb int
	failLiveness    bool
	failLivenessNb  int
}

func (s *MonitoringEndpoints) Liveness(l *log.Entry, w http.ResponseWriter, r *http.Request) {
	s.lock.Lock()
	defer s.lock.Unlock()

	if s.failLiveness {
		w.WriteHeader(s.config.LivenessStatusError)
	} else if s.failLivenessNb > 0 {
		s.failLivenessNb--
		w.WriteHeader(s.config.LivenessStatusError)
	}
	w.WriteHeader(s.config.LivenessStatusOk)
}

func (s *MonitoringEndpoints) Readiness(l *log.Entry, w http.ResponseWriter, r *http.Request) {
	s.lock.Lock()
	defer s.lock.Unlock()

	if s.failReadiness {
		w.WriteHeader(s.config.ReadinessStatusError)
	} else if s.failReadinessNb > 0 {
		s.failReadinessNb--
		w.WriteHeader(s.config.ReadinessStatusError)
	}
	w.WriteHeader(s.config.ReadinessStatusOk)
}

func (s *MonitoringEndpoints) FailReadiness(l *log.Entry, w http.ResponseWriter, r *http.Request) {
	failReadinessString := r.URL.Query().Get(monitoringQueryParamFail)
	if len(failReadinessString) == 0 {
		errString := fmt.Sprintf("missing %q query param", monitoringQueryParamFail)
		l.Errorf(errString)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(errString))
		return
	}
	if failReadiness, err := strconv.ParseBool(failReadinessString); err != nil {
		l.WithError(err).Errorf("failed to parse %q query param", monitoringQueryParamFail)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "failed to parse %q query param: %s", monitoringQueryParamFail, err.Error())
	} else {
		s.lock.Lock()
		s.failReadiness = failReadiness
		s.lock.Unlock()
		w.WriteHeader(http.StatusOK)
	}
}
