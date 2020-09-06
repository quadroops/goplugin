package goplugin

import (
	"fmt"

	"github.com/quadroops/goplugin/pkg/caller"
	"github.com/quadroops/goplugin/pkg/caller/driver"
)

// ProtocolOption used to configure rest or grpc options
// You have to choose between rest or grpc, if your plugin
// using rest, than ignore grpc option, and vice versa, but you
// can't ignore them both
type ProtocolOption struct {
	RestAddr    string
	RestTimeout int
	GrpcClient  driver.MakeClient
}

// BuildProtocol a helper to check and also generate a default
// protocol should be used like a REST or GRPC
func BuildProtocol(opt *ProtocolOption) caller.Builder {
	if opt == nil {
		return nil
	}

	// if both of configurations is not defined then return nil
	// no need to continue
	if opt.GrpcClient == nil && opt.RestAddr == "" {
		return nil
	}

	return func(commType string, port int) caller.Caller {
		switch commType {
		case "rest":
			addr := fmt.Sprintf("%s:%d", opt.RestAddr, port)
			return driver.NewREST(addr, &driver.RESTOption{
				Timeout: opt.RestTimeout,
			})
		case "grpc":
			return driver.NewGRPC(opt.GrpcClient)
		}

		return nil
	}
}
