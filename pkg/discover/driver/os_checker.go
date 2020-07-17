package driver

import (
	"fmt"
	"os"

	"github.com/quadroops/goplugin/pkg/discover"
)

const (
	// OSEnvName is an os variable name
	OSEnvName = "GOPLUGIN_DIR"
)

type osChecker struct {}

// NewOsChecker used to create new instance os checker
func NewOsChecker() discover.Checker {
	return new(osChecker) 
}

func (c *osChecker) Check() string {
	val := os.Getenv(OSEnvName)
	if val != "" {
		_, err := os.Stat(val)
		if os.IsNotExist(err) {
			panic(fmt.Sprintf("Path is not exist: %s", val))
		}

		return val
	}

	return ""
}