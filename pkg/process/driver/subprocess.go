package driver

import (
	"context"
	"fmt"
	"log"
	"os/exec"
	"syscall"
	"time"

	"github.com/quadroops/goplugin/pkg/errs"
	"github.com/quadroops/goplugin/pkg/process"
	"github.com/quadroops/goplugin/internal/utils"
)

type runner struct {}

// NewSubProcess used to create new instance that implement Runner
func NewSubProcess() process.Runner {
	return &runner{}
}

func (r *runner) Run(toWait int, name, command string, args ...string) (<-chan process.Plugin, error) {
	var stdout, stderr utils.Buffer 
	ctx, cancel := context.WithCancel(context.Background())

	cmd := exec.CommandContext(ctx, command, args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Start()
	if err != nil {
		cancel() // manually cancel the context and kill the process
		return nil, fmt.Errorf("%q: %w", err.Error(), errs.ErrPluginCannotStart) 
	}

	if toWait > 0 {
		// waiting the process
		time.Sleep(time.Duration(toWait) * time.Second)
	}

	ch := make(chan process.Plugin)
	go func() {
		plugin := process.Plugin{
			Kill: cancel,
			ID: process.ID(cmd.Process.Pid),
			Name: name,
			Stderr: &stderr,
			Stdout: &stdout,
		}

		ch <- plugin
		close(ch)
	}()

	go func() {
		// for now we doesn't need to handle an error from process
		// just log the error message
		err = cmd.Wait()
		if err != nil {
			log.Printf("Error wait: %v", err)
		}
	}()

	return ch, nil
}