package driver_test

import (
	"testing"

	"github.com/quadroops/goplugin/pkg/discover/driver"
	"github.com/stretchr/testify/assert"
)

const (
	fileToTest = "./tmp/test"
)

func TestFileReaderSucces(t *testing.T) {
	fr := driver.NewFileReader()
	b, err := fr.ReadFile(fileToTest)
	assert.NoError(t, err)
	assert.Equal(t, string(b), "testing")
}

func TestFileReaderFileNotFound(t *testing.T) {
	fr := driver.NewFileReader()
	_, err := fr.ReadFile("Unknown file")
	assert.Error(t, err)
}