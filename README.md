# Child process manager for go

This is a simple library to ensure all the child processes are killed if the main process is killed.

## Usage

Import the module

```
go get github.com/AgustinSRG/go-child-process-manager
```

Example usage:

```go
package main

import (
	"os/exec"

	// Import the module
	child_process_manager "github.com/AgustinSRG/go-child-process-manager"
)

func main() {
	// The child process manager must be initialized before any child processes are created
	err := child_process_manager.InitalizeChildProcessManager()
	if err != nil {
		panic(err)
	}
	// Call DisposeChildProcessManager() just before exiting the main process
	defer child_process_manager.DisposeChildProcessManager()

	// Create a command (this is an example)
	cmd := exec.Command("ffmpeg", "-i", "input.mp4", "output.webm")

	// Configure the command to be killed when the main process dies
	err = child_process_manager.ConfigureCommand(cmd)
	if err != nil {
		panic(err)
	}

	// Start the process
	err = cmd.Start()
	if err != nil {
		panic(err)
	}

	// Add process as a child process
	err = child_process_manager.AddChildProcess(cmd.Process)
	if err != nil {
		cmd.Process.Kill() // We must kill the process if this fails
		panic(err)
	}

	// Wait for the process to finish
	err = cmd.Wait()
	if err != nil {
		panic(err)
	}
}
```

## Testing

In order to test the code, first go to the `test` folder.

```sh
cd test
```

Then, build the test binary:

```sh
go get github.com/AgustinSRG/go-child-process-manager/test
go build
```

Run the test binary:

```sh
./test
```

After it finishes, it should print `SUCCESS | THE CHILD PROCESS IS DEAD`

If you want to test what happens without using this library, use the following:

```sh
./test --no-library
```

After it finishes, it should print `ERROR | THE CHILD PROCESS WAS ALIVE`
