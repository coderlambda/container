package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"syscall"
)

// #define _GNU_SOURCE
// #include <unistd.h>
// #include <sched.h>
import "C"

func runCmd(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println(name)
		fmt.Println(args)
		fmt.Println(err.Error())
		return err
	}
	return nil
}

func main() {
	pipeReader, pipeWriter, err := os.Pipe()

	unshareErr := C.unshare(syscall.CLONE_NEWUTS | syscall.CLONE_NEWIPC | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS)
	if unshareErr != 0 {
		fmt.Println("unshare error", unshareErr)
	}

	childPid, err := C.fork()
	if err != nil {
		fmt.Println(err.Error())
	}

	if childPid == 0 {
		unshareErr := C.unshare(syscall.CLONE_NEWNET)
		pipeWriter.Close()
		buffer := make([]byte, 4)
		pipeReader.Read(buffer)

		if unshareErr != 0 {
			fmt.Println("unshare error", unshareErr)
		}

		err := processInChild()
		if err != nil {
			fmt.Println(err.Error())
		}
	} else {
		f, _ := os.Open("/proc/self/ns/pid")
		C.setns(C.int(f.Fd()), syscall.CLONE_NEWPID)
		runCmd("ip", "link", "add", "veth0", "type", "veth", "peer", "name", "veth1")
		runCmd("ip", "link", "set", "veth1", "netns", strconv.Itoa(int(childPid)))
		runCmd("ip", "link", "set", "veth0", "up")
		runCmd("ip", "addr", "add", "169.254.1.1/30", "dev", "veth0")
		runCmd("iptables", "-t", "nat", "-A", "POSTROUTING", "-m", "veth0", "-j", "MASQUERADE")
		pipeWriter.Write([]byte("0"))

		var ws syscall.WaitStatus
		syscall.Wait4(int(childPid), &ws, 0, nil)
	}
}

func processInChild() error {
	C.sethostname(C.CString("in-namespace"), 12)
	syscall.Mount("proc", "/proc", "proc", 0, "")
	runCmd("ip", "link", "set", "lo", "up")
	runCmd("ip", "link", "set", "veth1", "up")
	runCmd("ip", "addr", "add", "169.254.1.2/30", "dev", "veth1")
	
	cmd := "bash"
	cmdFullPath, lookErr := exec.LookPath(cmd)
	if lookErr != nil {
		panic(lookErr)
	}

	args := []string{cmd}

	env := os.Environ()

	execErr := syscall.Exec(cmdFullPath, args, env)
	if execErr != nil {
		panic(execErr)
	}

	return nil
}
