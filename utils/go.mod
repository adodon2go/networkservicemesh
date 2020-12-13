module cisco-app-networking.github.io/networkservicemesh/utils

go 1.13

require (
	github.com/caddyserver/caddy v1.0.5
	github.com/networkservicemesh/networkservicemesh/controlplane/api v0.0.0-00010101000000-000000000000
	github.com/onsi/gomega v1.7.0
	github.com/pkg/errors v0.8.1
	github.com/sirupsen/logrus v1.4.2
	github.com/vishvananda/netns v0.0.0-20190625233234-7109fa855b0f
)

replace github.com/census-instrumentation/opencensus-proto v0.1.0-0.20181214143942-ba49f56771b8 => github.com/census-instrumentation/opencensus-proto v0.0.3-0.20181214143942-ba49f56771b8

replace (
	github.com/networkservicemesh/networkservicemesh/controlplane/api => cisco-app-networking.github.io/networkservicemesh/controlplane/api latest
	github.com/networkservicemesh/networkservicemesh/utils => ./
)
