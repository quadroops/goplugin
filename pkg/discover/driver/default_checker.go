package driver

import (
	"fmt"
	"os"

	"github.com/quadroops/goplugin/pkg/discover"
	"github.com/mitchellh/go-homedir"
)

const (
	// DefaultCheckerFilePath used as default config file name
	DefaultCheckerFilePath = ".goplugin/config.toml"
)

type defaultChecker struct {}

// NewDefaultChecker used to create new instance of default checker
func NewDefaultChecker() discover.Checker {
	return new(defaultChecker)
}

func (dc *defaultChecker) Check() string {
	dir, err := homedir.Dir()
	if err != nil {
		panic(err)
	}

	filepath := fmt.Sprintf("%s/%s", dir, DefaultCheckerFilePath)
	_, err = os.Stat(filepath)
	if os.IsNotExist(err) {
		panic(fmt.Sprintf("Configuration file not found at: %s", filepath))
	}

	return filepath
}