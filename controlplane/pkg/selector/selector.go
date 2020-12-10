package selector

import (
	"github.com/adodon2go/networkservicemesh/controlplane/api/connection"
	"github.com/adodon2go/networkservicemesh/controlplane/api/registry"
)

type Selector interface {
	SelectEndpoint(requestConnection *connection.Connection, ns *registry.NetworkService, networkServiceEndpoints []*registry.NetworkServiceEndpoint) *registry.NetworkServiceEndpoint
}
