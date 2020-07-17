package driver_test

import (
	"os"
	"testing"

	"github.com/quadroops/goplugin/pkg/discover/driver"
	"github.com/stretchr/testify/assert"
)

func TestEnvVarExist(t *testing.T) {
	os.Setenv(driver.OSEnvName, fileToTest)
	osChecker := driver.NewOsChecker()
	assert.NotPanics(t, func () {
		filepath := osChecker.Check() 
		assert.Equal(t, filepath, fileToTest)
	})
}

func TestEnvVarNotExist(t *testing.T) {
	os.Setenv(driver.OSEnvName, "")
	osChecker := driver.NewOsChecker()
	assert.NotPanics(t, func () {
		filepath := osChecker.Check() 
		assert.Empty(t, filepath)
	})
}

func TestPanic(t *testing.T) {
	os.Setenv(driver.OSEnvName, "test")
	osChecker := driver.NewOsChecker()
	assert.Panics(t, func() {
		osChecker.Check() 
	})
}