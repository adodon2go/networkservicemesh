package nsmmonitor

import (
	"cisco-app-networking.github.io/networkservicemesh/utils"
	"cisco-app-networking.github.io/networkservicemesh/utils/caddyfile"
	"cisco-app-networking.github.io/networkservicemesh/utils/dnsconfig"

	"github.com/sirupsen/logrus"

	"cisco-app-networking.github.io/networkservicemesh/controlplane/api/connection"
)

//nsmDNSMonitorHandler implements Handler interface for handling dnsConfigs
type nsmDNSMonitorHandler struct {
	EmptyNSMMonitorHandler
	manager  dnsconfig.Manager
	reloadOp utils.Operation
	path     string
}

func (m *nsmDNSMonitorHandler) Updated(old, new *connection.Connection) {
	logrus.Infof("Deleting config with id %v", old.Id)
	m.manager.Delete(old.Id)
	logrus.Infof("Adding config with id %v", new.Id)
	m.manager.Store(new.Id, new.GetContext().GetDnsContext().GetConfigs()...)
	m.reloadOp.Run()
}

//NewNsmDNSMonitorHandler creates new DNS monitor handler
func NewNsmDNSMonitorHandler() Handler {
	p := caddyfile.Path()
	mgr, err := dnsconfig.NewManagerFromCaddyfile(p)
	if err != nil {
		logrus.Fatalf("An error during parse corefile: %v", err)
	}
	m := &nsmDNSMonitorHandler{
		manager: mgr,
		path:    p,
	}
	m.reloadOp = utils.NewSingleAsyncOperation(func() {
		err := m.manager.Caddyfile(m.path).Save()
		if err != nil {
			logrus.Error(err)
		}
	})
	return m
}

func (m *nsmDNSMonitorHandler) Connected(conns map[string]*connection.Connection) {
	for _, conn := range conns {
		logrus.Info(conn.Context.DnsContext)
		err := m.manager.Caddyfile(m.path).Save()
		if err != nil {
			logrus.Error(err)
		}
		m.manager.Store(conn.Id, conn.GetContext().GetDnsContext().GetConfigs()...)
	}
	m.reloadOp.Run()
}

func (m *nsmDNSMonitorHandler) Closed(conn *connection.Connection) {
	logrus.Infof("Deleting config with id %v", conn.Id)
	m.manager.Delete(conn.Id)
	m.reloadOp.Run()
}
