// Testing program for go-child-process-manager

package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	child_process_manager "github.com/AgustinSRG/go-child-process-manager"
)

// Test entry point
func main() {
	if len(os.Args) < 1 {
		return
	}

	useLibrary := true
	isIntermediary := false
	isChild := false

	for i := 1; i < len(os.Args); i++ {
		if os.Args[i] == "--intermediary" {
			isIntermediary = true
		} else if os.Args[i] == "--child" {
			isChild = true
		} else if os.Args[i] == "--no-library" {
			useLibrary = false
		} else {
			fmt.Println("Unrecognized option: " + os.Args[i])
			os.Exit(1)
		}
	}

	if isIntermediary {
		RunIntermediaryProcess(useLibrary)
	} else if isChild {
		RunChildProcess()
	} else {
		RunMainProcess(useLibrary)
	}
}

func RunMainProcess(useLibrary bool) {
	fmt.Println("Running intermediary...")
	cmd := exec.Command(os.Args[0], "--intermediary")

	if !useLibrary {
		cmd.Args = append(cmd.Args, "--no-library")
	}

	// Create a pipe to read StdErr
	pipe, err := cmd.StderrPipe()

	if err != nil {
		panic(err)
	}

	if err = cmd.Start(); err != nil {
		panic(err)
	}

	fmt.Println("Intermediary has PID = " + fmt.Sprint(cmd.Process.Pid))

	fmt.Println("Waiting for intermediary to print child process PID...")

	reader := bufio.NewReader(pipe)

	line, err := reader.ReadString('\n')

	line = strings.ReplaceAll(line, "\n", "")

	pid, err := strconv.ParseInt(line, 10, 64)

	if err != nil {
		panic(err)
	}

	fmt.Println("Child process has PID = " + fmt.Sprint(pid))

	fmt.Println("Killing Intermediary...")

	cmd.Process.Kill()

	fmt.Println("Waiting 5 seconds...")

	time.Sleep(5 * time.Second)

	fmt.Println("Checking if child process is alive...")

	childP, err := os.FindProcess(int(pid))

	if err != nil {
		fmt.Println("SUCESS | THE CHILD PROCESS IS DEAD")
		return
	}

	err = childP.Kill()

	if err == nil {
		fmt.Println("ERROR | THE CHILD PROCESS WAS ALIVE")
		os.Exit(1)
	} else {
		fmt.Println("SUCESS | THE CHILD PROCESS IS DEAD")
	}
}

func RunChildProcess() {
	for {
		time.Sleep(1 * time.Second)
	}
}

func RunIntermediaryProcess(useLibrary bool) {
	if useLibrary {
		err := child_process_manager.InitalizeChildProcessManager()
		if err != nil {
			panic(err)
		}
		defer child_process_manager.DisposeChildProcessManager()
	}

	cmd := exec.Command(os.Args[0], "--child")

	if !useLibrary {
		cmd.Args = append(cmd.Args, "--no-library")
	}

	if useLibrary {
		err := child_process_manager.ConfigureCommand(cmd)
		if err != nil {
			panic(err)
		}
	}

	if err := cmd.Start(); err != nil {
		panic(err)
	}

	if useLibrary {
		if err := child_process_manager.AddChildProcess(cmd.Process); err != nil {
			panic(err)
		}
	}

	fmt.Fprint(os.Stderr, fmt.Sprint(cmd.Process.Pid)+"\n")

	if err := cmd.Wait(); err != nil {
		panic(err)
	}
}
