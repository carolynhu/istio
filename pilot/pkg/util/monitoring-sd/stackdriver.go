package monitoring_sd

import (
	"strings"

	"go.opencensus.io/stats/view"
)

const StackdriverPrefix = "stackdriver_"

func UseStackdiverStandardMetrics() bool {
	return true
	//return features.EnableStackdriverMetrics.Get()
}

func GetStackDriverMetricType(view *view.View) string {
	if UseStackdiverStandardMetrics() {
		if strings.HasPrefix(view.Name, StackdriverPrefix) {
			return "istio.io/control/" + strings.TrimPrefix(view.Name, StackdriverPrefix)
		}
	}
	return "custom.googleapis.com/istio_control/" + view.Name
}
