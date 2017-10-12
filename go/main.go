// +build linux

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"syscall"
)

const (
	CHROOT_PATH = "/tmp/smoker/rootfs"
	CG_NAME     = "smoker"
	CG_ROOT     = "/sys/fs/cgroup"
	ALPINE_URI  = "http://dl-cdn.alpinelinux.org/alpine/v3.6/releases/x86_64/alpine-minirootfs-3.6.2-x86_64.tar.gz"
)

func main() {
	switch os.Args[1] {
	case "run":
		runSelf(os.Args[2:])
	case "rmi":
		err := removeRootfs()
		if err != nil {
			log.Fatal(err)
		}
	case "-":
		runChild(os.Args[2:])
	case "pull":
		err := pullUnpack()
		if err != nil {
			log.Fatal("pulling image: ", err)
		}
		log.Println("image pulled and installed")
	}
}

// runSelf relaunches self
// process to make certain syscalls
// act only on child process.
func runSelf(args []string) {
	args = append([]string{"-"}, args...)
	cmd := exec.Command("/proc/self/exe", args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS,
		// /proc/self/exe will not be found
		// if we chroot here
		//Chroot: CHROOT_PATH,
	}
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}

func runChild(args []string) {
	if err := setCGroups(&args); err != nil {
		log.Fatal("cgroups: ", err)
	}
	err := os.MkdirAll(CHROOT_PATH, 0755)
	if err != nil {
		log.Fatal("creating chroot path: ", err)
	}
	if err = syscall.Chroot(CHROOT_PATH); err != nil {
		log.Fatal("chroot: ", err)
	}
	// necessary or the launched
	// command will have an undefined PWD
	if err := os.Chdir("/"); err != nil {
		log.Fatal("changing dir: ", err)
	}
	if err := syscall.Mount("proc", "/proc", "proc", 0, ""); err != nil {
		log.Fatal("mounting: ", err)
	}
	if err := syscall.Sethostname([]byte("my-container")); err != nil {
		log.Fatal("setting hostname: ", err)
	}
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
	if err := syscall.Unmount("/proc", 0); err != nil {
		log.Fatal("umount: ", err)
	}
}

// setCGroup configures control group
// to limit reasource consumption of container
func setCGroups(args *[]string) error {
	// avoid PID exhaustion
	// kernel/Documentation/cgroups/pids.txt
	if err := createCGroup("pids", "10"); err != nil {
		return err
	}
	if (*args)[0] == "-m" {
		if err := createCGroup("memory", (*args)[1]); err != nil {
			return err
		}
		*args = (*args)[2:]
	}
	return nil
}

// createCGroup creates cgroups for
// specific resource (pids, memory)
// and sets desired limit
func createCGroup(resource, limit string) error {
	resourceMap := map[string]string{
		"pids":   "/pids.max",
		"memory": "/memory.limit_in_bytes",
	}
	resourcePath := fmt.Sprintf("%s/%s/%s", CG_ROOT, resource, CG_NAME)
	err := os.MkdirAll(resourcePath, 0755)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(resourcePath+resourceMap[resource], []byte(limit), 0644)
	if err != nil {
		return err
	}
	if resource == "memory" {
		// you shall not swap!
		err = ioutil.WriteFile(resourcePath+"/memory.swappiness", []byte("0"), 0644)
		if err != nil {
			return err
		}
	}
	// TODO implement
	// https://access.redhat.com/documentation/en-us/red_hat_enterprise_linux/6/html/resource_management_guide/sec-common_tunable_parameters#ex-automatically_removing_empty_cgroups
	// automatically remove cgroup after
	// process exits
	//	err = ioutil.WriteFile(resourcePath+"/notify_on_release", []byte("1"), 0644)
	//	if err != nil {
	//		return err
	//	}
	//  ...
	p := strconv.Itoa(os.Getpid())
	// Are /tasks and /cgrorp.procs really the same thing?
	err = ioutil.WriteFile(resourcePath+"/tasks", []byte(p), 0644)
	if err != nil {
		return err
	}
	return nil
}

func pullUnpack() error {
	res, err := http.Get(ALPINE_URI)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	err = os.MkdirAll(CHROOT_PATH, 0755)
	if err != nil {
		return err
	}
	// lazily spawn a tar process
	// to avoid the pitfalls of tar package
	cmd := exec.Command("tar", "-xzf", "-", "-C", CHROOT_PATH)
	cmd.Stdin = res.Body
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func removeRootfs() error {
	return os.RemoveAll(CHROOT_PATH)
}
