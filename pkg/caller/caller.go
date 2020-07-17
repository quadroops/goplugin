package caller

import (
	"github.com/quadroops/goplugin/pkg/host"
)

// New used to create plugin's instance
func New(meta *host.Registry, transporter Caller) *Plugin {
	return &Plugin{meta, transporter}
}

// Ping used to send ping request to plugin
func (p *Plugin) Ping() (string, error) {
	return p.transporter.Ping()
}

// Exec used to send exec request to plugin
func (p *Plugin) Exec(cmdName string, payload []byte) ([]byte, error) {
	return p.transporter.Exec(cmdName, payload)
}