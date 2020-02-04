package bootstrap

import (
	"strings"

	"go.opencensus.io/stats/view"
	"istio.io/istio/pilot/pkg/features"
	"istio.io/pkg/env"
)

const StackdriverPrefix = "stackdriver_"

func UseStackdiverStandardMetrics() env.BoolVar {
	// this should be based on some env var or flag
	return features.EnableStackdriverMetrics
}

func GetStackDriverMetricType(view *view.View) string {
	if UseStackdiverStandardMetrics().DefaultValue == "true" {
		if strings.HasPrefix(view.Name, StackdriverPrefix) {
			return "istio.io/control/" + strings.TrimPrefix(view.Name, StackdriverPrefix)
		}
	}
	return "custom.googleapis.com/istio_control" + view.Name
}
