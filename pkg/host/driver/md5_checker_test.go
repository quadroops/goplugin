package driver_test

import (
	"testing"

	"github.com/quadroops/goplugin/pkg/host/driver"
	"github.com/stretchr/testify/assert"
)

func TestGetMD5Success(t *testing.T) {
	drv := driver.NewMd5Check()
	sum, err := drv.Parse("./tmp/test")
	assert.NoError(t, err)
	assert.NotEmpty(t, sum)
}

func TestGetMD5FileNotExist(t *testing.T) {
	drv := driver.NewMd5Check()
	sum, err := drv.Parse("./tmp/test2")
	assert.Error(t, err)
	assert.Empty(t, sum)
}