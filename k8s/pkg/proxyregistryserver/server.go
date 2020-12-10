package proxyregistryserver

import (
	"context"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"github.com/adodon2go/networkservicemesh/pkg/tools"

	"github.com/adodon2go/networkservicemesh/controlplane/api/clusterinfo"
	"github.com/adodon2go/networkservicemesh/controlplane/api/registry"
	nsmClientset "github.com/adodon2go/networkservicemesh/k8s/pkg/networkservice/clientset/versioned"
	"github.com/adodon2go/networkservicemesh/k8s/pkg/registryserver"
)

// New starts proxy Network Service Discovery Server and Cluster Info Server
func New(clientset *nsmClientset.Clientset, clusterInfoService clusterinfo.ClusterInfoServer) *grpc.Server {
	server := tools.NewServer(context.Background())
	cache := registryserver.NewRegistryCache(clientset, &registryserver.ResourceFilterConfig{})
	discovery := newDiscoveryService(cache, clusterInfoService)
	nseRegistry := newNseRegistryService(clusterInfoService)

	registry.RegisterNetworkServiceDiscoveryServer(server, discovery)
	registry.RegisterNetworkServiceRegistryServer(server, nseRegistry)
	clusterinfo.RegisterClusterInfoServer(server, clusterInfoService)

	if err := cache.Start(); err != nil {
		logrus.Error(err)
	}
	logrus.Info("RegistryCache started")

	return server
}
