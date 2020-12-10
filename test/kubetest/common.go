package kubetest

import (
	"github.com/sirupsen/logrus"

	"github.com/adodon2go/networkservicemesh/test/kubetest/pods"
)

// DefaultForwarderVariables - Default variables for forwarder deployment
func DefaultForwarderVariables(plane string) map[string]string {
	if plane == pods.EnvForwardingPlaneDefault {
		return DefaultPlaneVariablesVPP()
	} else if plane == pods.EnvForwardingPlaneKernel {
		return DefaultPlaneVariablesKernel()
	}
	logrus.Error("Forwarding plane error: Unknown forwarder")
	return nil
}
