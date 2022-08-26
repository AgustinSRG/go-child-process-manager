//go:build windows

// Original source of this aproach: https://gist.github.com/hallazzang/76f3970bfc949831808bbebc8ca15209
// Credits for this code to @hallazzang

package child_process_manager

import (
	"os"
	"os/exec"
	"sync"
	"unsafe"

	"golang.org/x/sys/windows"
)

// We use this struct to retreive process handle(which is unexported)
// from os.Process using unsafe operation.
type process struct {
	Pid    int
	Handle uintptr
}

// Child process manager for windows
type childProcessManager struct {
	exit_group windows.Handle // Windows exit group
	mu         *sync.Mutex    // Mutex for multi-threading access
}

var (
	child_process_manager *childProcessManager = nil
)

func createChildProcessManager() (*childProcessManager, error) {
	handle, err := windows.CreateJobObject(nil, nil)
	if err != nil {
		return nil, err
	}

	info := windows.JOBOBJECT_EXTENDED_LIMIT_INFORMATION{
		BasicLimitInformation: windows.JOBOBJECT_BASIC_LIMIT_INFORMATION{
			LimitFlags: windows.JOB_OBJECT_LIMIT_KILL_ON_JOB_CLOSE,
		},
	}
	if _, err := windows.SetInformationJobObject(
		handle,
		windows.JobObjectExtendedLimitInformation,
		uintptr(unsafe.Pointer(&info)),
		uint32(unsafe.Sizeof(info))); err != nil {
		return nil, err
	}

	return &childProcessManager{
		exit_group: handle,
		mu:         &sync.Mutex{},
	}, nil
}

// Initializes child process manager
func InitalizeChildProcessManager() error {
	pm, err := createChildProcessManager()

	if err != nil {
		return err
	}

	child_process_manager = pm

	return nil
}

// Configures the command to be a child process, use this before you call Start()
// cmd - Command reference
func ConfigureCommand(cmd *exec.Cmd) error {
	return nil
}

// Adds a child process, so it is killed if the main process dies
// p - Child process
func AddChildProcess(p *os.Process) error {
	child_process_manager.mu.Lock()
	defer child_process_manager.mu.Unlock()

	return windows.AssignProcessToJobObject(
		windows.Handle(child_process_manager.exit_group),
		windows.Handle((*process)(unsafe.Pointer(p)).Handle))
}

// Disposes resources created for the process manager
// Call this before the main process ends
func DisposeChildProcessManager() error {
	child_process_manager.mu.Lock()
	defer child_process_manager.mu.Unlock()

	return windows.CloseHandle(windows.Handle(child_process_manager.exit_group))
}
