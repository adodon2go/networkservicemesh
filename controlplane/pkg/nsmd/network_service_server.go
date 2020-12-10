package nsmd

import (
	"github.com/adodon2go/networkservicemesh/controlplane/api/networkservice"
	"github.com/adodon2go/networkservicemesh/controlplane/pkg/api/nsm"
	"github.com/adodon2go/networkservicemesh/controlplane/pkg/common"
	"github.com/adodon2go/networkservicemesh/controlplane/pkg/local"
	"github.com/adodon2go/networkservicemesh/controlplane/pkg/model"
)

// NewNetworkServiceServer - construct a local network service chain
func NewNetworkServiceServer(model model.Model, ws *Workspace,
	nsmManager nsm.NetworkServiceManager) networkservice.NetworkServiceServer {
	return common.NewCompositeService("Local",
		common.NewRequestValidator(),
		common.NewMonitorService(ws.MonitorConnectionServer()),
		local.NewWorkspaceService(ws.Name()),
		local.NewConnectionService(model),
		local.NewForwarderService(model, nsmManager.ServiceRegistry()),
		local.NewEndpointSelectorService(nsmManager.NseManager()),
		common.NewExcludedPrefixesService(),
		local.NewEndpointService(nsmManager.NseManager(), nsmManager.GetHealProperties(), nsmManager.Model()),
		common.NewCrossConnectService(),
	)
}
