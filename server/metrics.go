package server

import (
	"github.com/cyberdelia/go-metrics-graphite"
	"github.com/rcrowley/go-metrics"
)

func startMetricsProcess(server *Server) {
	if len(server.graphiteConfigs) > 0 {
		log.Debug("Starting sending runtime metrics to graphite, config %+v", server.graphiteConfigs)
		metrics.RegisterRuntimeMemStats(server.graphiteConfigs[0].Registry)
		go metrics.CaptureRuntimeMemStats(server.graphiteConfigs[0].Registry, server.graphiteConfigs[0].FlushInterval)
		for _, config := range server.graphiteConfigs {
			go graphite.WithConfig(config)
		}

	}
}