package goplugin

import (
	"github.com/quadroops/goplugin/pkg/caller"
	"github.com/quadroops/goplugin/pkg/caller/driver"
)

// BuildProtocol a helper to check and also generate a default
// protocol should be used like a REST or GRPC
func BuildProtocol(opt *ProtocolOption) caller.Builder {
	if opt == nil {
		return nil
	}

	// if both of configurations is not defined then return nil
	// no need to continue
	if opt.GRPCOpts == nil && opt.RESTOpts == nil {
		return nil
	}

	return func(commType string, port int) caller.Caller {
		switch commType {
		case "rest":
			return driver.NewREST(opt.RESTOpts)
		case "grpc":
			return driver.NewGRPC(opt.GRPCOpts)
		}

		return nil
	}
}
