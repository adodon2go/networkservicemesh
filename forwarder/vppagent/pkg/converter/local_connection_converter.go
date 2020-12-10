package converter

import (
	"github.com/adodon2go/networkservicemesh/controlplane/api/connection"
	"github.com/adodon2go/networkservicemesh/controlplane/api/connection/mechanisms/kernel"
	"github.com/adodon2go/networkservicemesh/controlplane/api/connection/mechanisms/memif"
)

type LocalConnectionConverter struct {
	*connection.Connection
	name         string
	ipAddressKey string
}

func NewLocalConnectionConverter(c *connection.Connection, conversionParameters *ConnectionConversionParameters) Converter {
	if c.GetMechanism().GetType() == kernel.MECHANISM {
		return NewKernelConnectionConverter(c, conversionParameters)
	}
	if c.GetMechanism().GetType() == memif.MECHANISM {
		return NewMemifInterfaceConverter(c, conversionParameters)
	}
	return nil
}
