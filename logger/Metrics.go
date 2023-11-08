package logger

import "github.com/smira/go-statsd"

var metricsClient *statsd.Client

func InitMetrics() *statsd.Client {
	client := statsd.NewClient("localhost:8125",
		statsd.MaxPacketSize(1400),
		statsd.MetricPrefix("web."))
	metricsClient = client
	return metricsClient
}

func GetMetricsClient() *statsd.Client {
	return metricsClient
}
