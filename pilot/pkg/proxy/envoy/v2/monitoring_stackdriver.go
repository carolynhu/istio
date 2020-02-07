package v2

import (
	monitoring_sd "istio.io/istio/pilot/pkg/util/monitoring-sd"
	"istio.io/pkg/monitoring"
)

var (
	configTag = monitoring.MustCreateLabel("config")

	configPushes = monitoring.NewSum(
		monitoring_sd.StackdriverPrefix+"config_pushes_count",
		"Istio control plane configuration pushes",
		monitoring.WithLabels(typeTag, configTag),
	)

	rdsConfigPushes = configPushes.With(typeTag.Value("RDS"))
	ldsConfigPushes = configPushes.With(typeTag.Value("LDS"))
	edsConfigPushes = configPushes.With(typeTag.Value("EDS"))
	cdsConfigPushes = configPushes.With(typeTag.Value("CDS"))

	rejectedConfigs = monitoring.NewGauge(
		monitoring_sd.StackdriverPrefix+"rejected_configs",
		"Istio control plane rejected configurations",
		monitoring.WithLabels(typeTag, errTag),
	)
)

func init() {
	if monitoring_sd.UseStackdiverStandardMetrics() {
		monitoring.MustRegister(
			configPushes,
			rejectedConfigs,
		)
	}
}
