package kubetest

import "cisco-app-networking.github.io/networkservicemesh/forwarder/pkg/common"

// DefaultPlaneVariablesKernel - Default variables for Kernel forwarding deployment
func DefaultPlaneVariablesKernel() map[string]string {
	return map[string]string{
		common.ForwarderMetricsEnabledKey: "false",
	}
}
