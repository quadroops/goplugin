package driver_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/quadroops/goplugin/pkg/discover/driver"
	"github.com/mitchellh/go-homedir"
	"github.com/stretchr/testify/assert"
)

func _getDir() string {
	dir, _ := homedir.Dir()
	dirpath := fmt.Sprintf("%s/.goplugin", dir) 
	return dirpath
}

func _getFilePath(dir string) string {
	return fmt.Sprintf("%s/config.toml", dir)
}

func _createHomeDirPath(dirpath string) {
	os.Mkdir(dirpath, os.ModeDir)	
}

func _removeHomeDirPath(dirpath string) {
	os.Remove(dirpath)
}

func _touch(filepath string) {
	f, _ := os.Create(filepath)
	f.Close()
}

func TestDefaultCheckerThrowPanic(t *testing.T) {
	_removeHomeDirPath(_getDir())
	checker := driver.NewDefaultChecker()
	assert.Panics(t, func() {
		checker.Check()
	})
}

func TestDefaultCheckerSuccess(t *testing.T) {
	_createHomeDirPath(_getDir())

	filepath := _getFilePath(_getDir())
	_touch(filepath)

	checker := driver.NewDefaultChecker()
	assert.NotPanics(t, func() {
		fileReturn := checker.Check()
		assert.NotEmpty(t, fileReturn)
		assert.Equal(t, filepath, fileReturn)
	})
}