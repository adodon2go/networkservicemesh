package converter_test

import (
	"fmt"
	"testing"

	"cisco-app-networking.github.io/networkservicemesh/controlplane/api/connection/mechanisms/kernel"

	"cisco-app-networking.github.io/networkservicemesh/forwarder/vppagent/pkg/converter"
)

func TestTempIf(t *testing.T) {
	tempIface := converter.TempIfName()
	fmt.Printf("tempIface: %s len(tempIface) %d\n", tempIface, len(tempIface))
	if len(tempIface) > kernel.LinuxIfMaxLength {
		t.Errorf("%s is longer than %d", tempIface, kernel.LinuxIfMaxLength)
	}
}
