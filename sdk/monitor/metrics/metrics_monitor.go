package metrics

import "github.com/adodon2go/networkservicemesh/controlplane/api/crossconnect"

type MetricsMonitor interface {
	HandleMetrics(statistics map[string]*crossconnect.Metrics)
}
