package caller

import "github.com/quadroops/goplugin/pkg/host"

var (
	// AllowedProtocols used as main supported protocols
	AllowedProtocols = []string{"rest", "grpc"}
)

// Builder used as a simple function to create Caller instance
type Builder func(commType string, port int) Caller

// Caller used as main communication interface
type Caller interface {
	Ping() (string, error)
	Exec(cmdName string, payload []byte) ([]byte, error)
}

// Plugin is single plugin instance used to store
// meta information and also caller activity
type Plugin struct {
	Meta        *host.Registry
	transporter Caller
}
