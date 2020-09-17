package driver_test

import (
	"errors"
	"log"
	"os"
	"syscall"
	"testing"

	"github.com/quadroops/goplugin/pkg/process/driver"
	"github.com/stretchr/testify/assert"
)

var (
	ErrFinished = errors.New("os: process already finished")
)

func checkProcessExist(id int) bool {
	p, err := os.FindProcess(id)
	if err != nil {
		log.Printf("Error find process: %v", err)
		return false
	}

	err = p.Signal(syscall.Signal(0))
	if err != nil {
		log.Printf("Error signal: %v | ProcessID: %v", err, id)
		if err.Error() != ErrFinished.Error() {
			log.Printf("Error unknown signal: %v", err)
		}
		return false
	}

	return true
}

func TestRunSubProcessSuccess(t *testing.T) {
	sub := driver.NewSubProcess()
	process, err := sub.Run(1, "test", "sleep", 5)
	assert.NoError(t, err)

	select {
	case plugin := <-process:
		assert.NotEmpty(t, plugin.ID)
		assert.NotEmpty(t, plugin.Name)
	default:
		// passed
	}
}

func TestRunSubProcessUnknownCommand(t *testing.T) {
	sub := driver.NewSubProcess()
	process, err := sub.Run(1, "test", "unkwon", 5)
	assert.Error(t, err)
	assert.Nil(t, process)
}
