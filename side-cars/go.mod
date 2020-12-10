module cisco-app-networking.github.io/networkservicemesh/side-cars

go 1.13

require (
	github.com/onsi/gomega v1.7.0
	github.com/pkg/errors v0.9.1
	github.com/sirupsen/logrus v1.4.2
	github.com/spiffe/spire/proto/spire v0.0.0-20200103215556-34b7e3785007
	google.golang.org/grpc v1.27.1
)

replace github.com/census-instrumentation/opencensus-proto v0.1.0-0.20181214143942-ba49f56771b8 => github.com/census-instrumentation/opencensus-proto v0.0.3-0.20181214143942-ba49f56771b8
