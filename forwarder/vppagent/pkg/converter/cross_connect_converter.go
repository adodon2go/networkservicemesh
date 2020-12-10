package converter

import (
	"path"

	vpp_l2 "go.ligato.io/vpp-agent/v3/proto/ligato/vpp/l2"

	"github.com/pkg/errors"

	"go.ligato.io/vpp-agent/v3/proto/ligato/configurator"
	"go.ligato.io/vpp-agent/v3/proto/ligato/vpp"

	"cisco-app-networking.github.io/networkservicemesh/controlplane/api/connection/mechanisms/common"
	"cisco-app-networking.github.io/networkservicemesh/controlplane/api/crossconnect"
)

const (
	srcPrefix = "SRC-"
	dstPrefix = "DST-"
)

type CrossConnectConverter struct {
	*crossconnect.CrossConnect
	conversionParameters *CrossConnectConversionParameters
}

func NewCrossConnectConverter(c *crossconnect.CrossConnect, conversionParameters *CrossConnectConversionParameters) *CrossConnectConverter {
	return &CrossConnectConverter{
		CrossConnect:         c,
		conversionParameters: conversionParameters,
	}
}

func (c *CrossConnectConverter) ToDataRequest(rv *configurator.Config, connect bool) (*configurator.Config, error) {
	if c == nil {
		return rv, errors.New("CrossConnectConverter cannot be nil")
	}
	if err := c.IsComplete(); err != nil {
		return rv, err
	}
	if rv == nil {
		rv = &configurator.Config{}
	}
	if rv.VppConfig == nil {
		rv.VppConfig = &vpp.ConfigData{}
	}

	srcName := srcPrefix + c.GetId()
	dstName := dstPrefix + c.GetId()
	mtu := c.calculateInterfaceMTU()

	if src := c.GetLocalSource(); src != nil {
		baseDir := path.Join(c.conversionParameters.BaseDir, src.GetMechanism().GetParameters()[common.Workspace])
		conversionParameters := &ConnectionConversionParameters{
			Name:      srcName,
			Terminate: false,
			Side:      SOURCE,
			BaseDir:   baseDir,
			MTU:       mtu,
		}
		var err error
		rv, err = NewLocalConnectionConverter(src, conversionParameters).ToDataRequest(rv, connect)
		if err != nil {
			return rv, errors.Wrapf(err, "Error Converting CrossConnect %v", c)
		}
	}

	if dst := c.GetLocalDestination(); dst != nil {
		baseDir := path.Join(c.conversionParameters.BaseDir, dst.GetMechanism().GetParameters()[common.Workspace])
		conversionParameters := &ConnectionConversionParameters{
			Name:      dstName,
			Terminate: false,
			Side:      DESTINATION,
			BaseDir:   baseDir,
			MTU:       mtu,
		}
		var err error
		rv, err = NewLocalConnectionConverter(dst, conversionParameters).ToDataRequest(rv, connect)
		if err != nil {
			return rv, errors.Wrapf(err, "Error Converting CrossConnect %v", c)
		}
	}

	rv, err := c.MechanismsToDataRequest(rv, connect)
	if err != nil {
		return rv, err
	}

	if len(rv.VppConfig.Interfaces) > 2 {
		return nil, errors.Errorf("created too many interfaces to cross connect, expected 2, got %d", len(rv.VppConfig.Interfaces))
	}

	// For connections mechanisms with xconnect required (For example SRv6 does not require xconnect)
	if len(rv.VppConfig.Interfaces) == 2 {
		ifaces := rv.VppConfig.Interfaces[len(rv.VppConfig.Interfaces)-2:]
		rv.VppConfig.XconnectPairs = append(rv.VppConfig.XconnectPairs, &vpp_l2.XConnectPair{
			ReceiveInterface:  ifaces[0].Name,
			TransmitInterface: ifaces[1].Name,
		})
		rv.VppConfig.XconnectPairs = append(rv.VppConfig.XconnectPairs, &vpp_l2.XConnectPair{
			ReceiveInterface:  ifaces[1].Name,
			TransmitInterface: ifaces[0].Name,
		})
	}

	return rv, nil
}

// MechanismsToDataRequest prepares data change with mechanisms parameters for vppagent
func (c *CrossConnectConverter) MechanismsToDataRequest(rv *configurator.Config, connect bool) (*configurator.Config, error) {
	if rv == nil {
		rv = &configurator.Config{}
	}
	if rv.VppConfig == nil {
		rv.VppConfig = &vpp.ConfigData{}
	}

	srcName := srcPrefix + c.GetId()
	dstName := dstPrefix + c.GetId()

	var err error
	if src := c.GetRemoteSource(); src != nil {
		rv, err = NewRemoteConnectionConverter(src, srcName, dstName, SOURCE).ToDataRequest(rv, connect)
		if err != nil {
			return rv, errors.Wrapf(err, "error Converting CrossConnect %v", c)
		}
	}

	if dst := c.GetRemoteDestination(); dst != nil {
		rv, err = NewRemoteConnectionConverter(dst, dstName, srcName, DESTINATION).ToDataRequest(rv, connect)
		if err != nil {
			return rv, errors.Wrapf(err, "error Converting CrossConnect %v", c)
		}
	}

	return rv, nil
}

// calculateInterfaceMTU returns the proper MTU to be applied on xconnect interfaces
func (c *CrossConnectConverter) calculateInterfaceMTU() uint32 {
	if c.conversionParameters.MTUOverride != 0 {
		return c.conversionParameters.MTUOverride // MTUOverride takes precedence if set
	}
	if c.conversionParameters.BaseMTU == 0 {
		return 0 // MTU 0 in vppagent API means undefined
	}
	// find the largest MTU overhead from both src/dst mechanisms and src/dst extra contexts
	srcOverhead, _ := c.Source.GetContext().GetMTUOverhead()
	dstOverhead, _ := c.Destination.GetContext().GetMTUOverhead()
	overheads := []uint32{
		c.conversionParameters.MechanismMTUOverhead,
		srcOverhead,
		dstOverhead,
	}
	maxOverhead := uint32(0)
	for _, o := range overheads {
		if o > maxOverhead {
			maxOverhead = o
		}
	}
	return c.conversionParameters.BaseMTU - maxOverhead
}
