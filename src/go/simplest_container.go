package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

// #define _GNU_SOURCE
// #include <unistd.h>
// #include <sched.h>
import "C"

func main() {

	unshareErr := C.unshare(syscall.CLONE_NEWUTS | syscall.CLONE_NEWIPC | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS | syscall.CLONE_NEWNET)

	if unshareErr != 0 {
		fmt.Println("unshare error", unshareErr)
	}

	pid, err := C.fork()
	if err != nil {
		fmt.Println(err.Error())
	}

	if pid == 0 {
		err := processInChild()
		if err != nil {
			fmt.Println(err.Error())
		}
	} else {
		var ws syscall.WaitStatus
		syscall.Wait4(int(pid), &ws, 0, nil)
	}
}

func processInChild() error {
	C.sethostname(C.CString("in-namespace"), 12)
	syscall.Mount("proc", "/proc", "proc", 0, "")
	cmd := "bash"
	binary, lookErr := exec.LookPath(cmd)
	if lookErr != nil {
		panic(lookErr)
	}

	args := []string{cmd}

	env := os.Environ()

	execErr := syscall.Exec(binary, args, env)
	if execErr != nil {
		panic(execErr)
	}

	return nil
}
