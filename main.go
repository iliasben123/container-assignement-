package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

func main() {
	if len(os.Args) < 2 {
		panic("Usage: run <command>")
	}
	switch os.Args[1] {
	case "run":
		run()
	case "child":
		child()
	default:
		panic("what")
	}
}

func run() {
	fmt.Printf("runing %v\n", os.Args[2:])
	cmd := exec.Command("/proc/self/exe", append([]string{"child"}, os.Args[2:]...)...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS | syscall.CLONE_NEWNET,
		
	}

	
	cgroup := "/sys/fs/cgroup/mon-docker"
	os.MkdirAll(cgroup, 0755)

	
	_ = os.WriteFile("/sys/fs/cgroup/cgroup.subtree_control", []byte("+memory"), 0700)

	
	must(os.WriteFile(cgroup+"/memory.max", []byte("50000000"), 0700))

	must(cmd.Start())

	
	must(os.WriteFile(cgroup+"/cgroup.procs", []byte(fmt.Sprintf("%d", cmd.Process.Pid)), 0700))

	must(cmd.Wait())
}

func child() {
	fmt.Printf("runing %v as pid %d\n", os.Args[2:], os.Getpid())

	must(syscall.Chroot("/home/utilisateur/mon-docker/rootfs"))
	must(os.Chdir("/"))

	
	flags := uintptr(syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV)
	must(syscall.Mount("proc", "proc", "proc", flags, ""))

	cmd := exec.Command(os.Args[2], os.Args[3:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	must(cmd.Run())

	
	must(syscall.Unmount("proc", 0))
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}