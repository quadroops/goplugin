package discover

import (
	"fmt"

	"github.com/quadroops/goplugin/pkg/errs"
)

// ConfigChecker used as main struct to explore config file path
type ConfigChecker struct {
	checkers []Checker
}

// NewConfigChecker used to create new instance of ConfigChecker
func NewConfigChecker(checkers ...Checker) *ConfigChecker {
	return &ConfigChecker{checkers: checkers}
}

// Explore used to explore configuration's file from os or default home dir
func (cc *ConfigChecker) Explore() (string, error) {
	if len(cc.checkers) < 1 {
		return "", fmt.Errorf("%w", errs.ErrDiscoverNoCheckers)
	}

	for _, checker := range cc.checkers {
		configPath := checker.Check()
		if configPath != "" && len(configPath) > 0 {
			return configPath, nil
		}
	}

	return "", fmt.Errorf("%w", errs.ErrConfigNotFound)
}
