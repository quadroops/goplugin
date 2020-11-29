package caller

import (
	"log"
	"time"

	"github.com/quadroops/goplugin/pkg/errs"
	"github.com/quadroops/goplugin/pkg/host"
)

var (
	// ignoredErrors used to put a list of error types
	// an error listed here, means, when that error happened
	// it will not trigger any retry mechanism
	ignoredErrors map[error]bool
)

func init() {
	ignoredErrors = make(map[error]bool)
	ignoredErrors[errs.ErrPluginPing] = true
	ignoredErrors[errs.ErrPluginExec] = true
	ignoredErrors[errs.ErrPluginCall] = true
}

// New used to create plugin's instance
func New(meta *host.Registry, transporter Caller, retryTimeout int) *Plugin {
	return &Plugin{meta, transporter, retryTimeout}
}

// Ping used to send ping request to plugin
func (p *Plugin) Ping() (string, error) {
	resp, err := p.transporter.Ping()
	if err != nil {
		if _, exist := ignoredErrors[err]; !exist {
			log.Printf("Retrying process error: %v", err)
			time.Sleep(time.Duration(p.retryTimeout) * time.Second)
			return p.Ping()
		}

		return "", err
	}

	return resp, nil
}

// Exec used to send exec request to plugin
func (p *Plugin) Exec(cmdName string, payload []byte) ([]byte, error) {
	resp, err := p.transporter.Exec(cmdName, payload)
	if err != nil {
		if _, exist := ignoredErrors[err]; !exist {
			log.Printf("Retrying process error: %v", err)
			time.Sleep(time.Duration(p.retryTimeout) * time.Second)
			return p.Exec(cmdName, payload)
		}

		return nil, err
	}

	return resp, nil
}
