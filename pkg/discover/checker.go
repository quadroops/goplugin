package discover

import (
	"fmt"

	"github.com/quadroops/goplugin/pkg/errs"
)

// ConfigChecker used as main struct to explore config file path
type ConfigChecker struct {
	osChecker      Checker
	defaultChecker Checker
}

// NewConfigChecker used to create new instance of ConfigChecker
func NewConfigChecker(osChecker, defaultChecker Checker) *ConfigChecker {
	return &ConfigChecker{osChecker: osChecker, defaultChecker: defaultChecker}
}

// Explore used to explore configuration's file from os or default home dir
func (cc *ConfigChecker) Explore() (string, error) {
	loadFromOS := cc.osChecker.Check()
	loadFromDefault := cc.defaultChecker.Check()

	if loadFromOS != "" {
		return loadFromOS, nil
	}

	if loadFromDefault != "" {
		return loadFromDefault, nil
	}

	return "", fmt.Errorf("%w", errs.ErrConfigNotFound)
}
