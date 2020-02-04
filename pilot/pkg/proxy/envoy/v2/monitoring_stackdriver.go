package v2

import (
	"istio.io/istio/pilot/pkg/bootstrap"
	"istio.io/pkg/monitoring"
)

var (
	configTag = monitoring.MustCreateLabel("config")

	configPushes = monitoring.NewSum(
		bootstrap.StackdriverPrefix+"config_pushes_count",
		"Istio control plane configuration pushes",
		monitoring.WithLabels(typeTag),
	)

	rdsConfigPushes = configPushes.With(typeTag.Value("RDS"))
	ldsConfigPushes = configPushes.With(typeTag.Value("LDS"))
	edsConfigPushes = configPushes.With(typeTag.Value("EDS"))
	cdsConfigPushes = configPushes.With(typeTag.Value("CDS"))

	rejectedConfigs = monitoring.NewGauge(
		bootstrap.StackdriverPrefix+"rejected_configs",
		"Istio control plane rejected configurations",
		monitoring.WithLabels(typeTag, errTag),
	)
)

func init() {
	if bootstrap.UseStackdiverStandardMetrics().DefaultValue == "true" {
		monitoring.MustRegister(
			configPushes,
			rejectedConfigs,
		)
	}
}
