// Copyright 2013 Albert P. Tobey. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"../src/lnxns"
	"fmt"
	"os"
	"syscall"
)

func main() {
	var root, cmd string
	var opts []string

	if len(os.Args) < 3 {
		panic("not enough arguments\n")
	}

	root = os.Args[1]
	cmd = os.Args[2]
	if len(os.Args) > 3 {
		opts = os.Args[3:len(os.Args)]
	}

	rs, err := os.Stat(root)
	if err != nil {
		if os.IsNotExist(err) {
			panic(fmt.Sprintf("'%s' does not exist, cannot chroot there!", root))
		} else {
			panic(fmt.Sprintf("Could not stat '%s': %s", root, err))
		}
	} else {
		if !rs.Mode().IsDir() {
			panic(fmt.Sprintf("'%s' is not a directory, cannot chroot there!", root))
		}
	}

	fmt.Printf("root: %s, cmd: %s, opts: %s\n", root, cmd, opts)

	err = os.Chdir(root)
	if err != nil {
		panic(fmt.Sprintf("chdir failed: %s", err))
	}

	err = syscall.Chroot(root)
	if err != nil {
		panic(fmt.Sprintf("chroot failed: %s", err))
	}

	// we're going to exec right away in the child, CLONE_VFORK will block the
	// parent from being scheduled until the child starts up, see clone(2)
	pid, err := lnxns.NsFork(lnxns.CLONE_VFORK)

	if err == syscall.EINVAL {
		panic("OS returned EINVAL. Make sure your kernel configuration includes all CONFIG_*_NS options.")
	} else if err != nil {
		panic(fmt.Sprintf("lnxns.NsFork() failed: %s", err))
	}

	if pid != 0 {
		proc, _ := os.FindProcess(pid)
		proc.Wait()
	} else {
		err = syscall.Exec(cmd, opts, os.Environ())
		if err != nil {
			panic(fmt.Sprintf("exec failed: %s", err))
		}
		panic("impossible")
	}
}

// vim: ts=4 sw=4 noet tw=120 softtabstop=4
