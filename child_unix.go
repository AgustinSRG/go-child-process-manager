//go:build unix

package child_process_manager

import (
	"os"
	"os/exec"
	"syscall"
)

// Initializes child process manager
func InitalizeChildProcessManager() error {
	return nil
}

// Configures the command to be a child process, use this before you call Start()
// cmd - Command reference
func ConfigureCommand(cmd *exec.Cmd) error {
	if cmd.SysProcAttr == nil {
		cmd.SysProcAttr = &syscall.SysProcAttr{}
	}
	cmd.SysProcAttr.Pdeathsig = syscall.SIGTERM
	return nil
}

// Adds a child process, so it is killed if the main process dies
// p - Child process
func AddChildProcess(p *os.Process) error {
	return nil
}

// Disposes resources created for the process manager
// Call this before the main process ends
func DisposeChildProcessManager() error {
	return nil
}
