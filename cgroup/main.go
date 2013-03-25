// Copyright 2013 Albert P. Tobey. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"../src/lnxns"
	"flag"
	"fmt"
	"os"
	"strings"
	"syscall"
)

type cgroupList []string
type envMap map[string]string

func (cgl *cgroupList) String() string {
	return fmt.Sprint(*cgl)
}

func (cgl *cgroupList) Set(value string) error {
	for _, cg := range strings.Split(value, ",") {
		*cgl = append(*cgl, cg)
	}
	return nil
}

func (env envMap) String() string {
	return fmt.Sprint(map[string]string(env))
}

func (env envMap) Set(value string) error {
	kv := strings.SplitN(value, "=", 2)
	env[kv[0]] = kv[1]
	return nil
}

var cgroupName string
var programFlag string
var argumentsFlag string
var envFlag envMap = make(envMap)
var cgRoot string

func init() {
	flag.StringVar(&cgroupName, "name", "lnxns", "name of the cgroup, must be a valid Linux directory name")
	flag.Var(&envFlag, "env", "key=value environment variables")
	flag.StringVar(&programFlag, "program", "", "the program to run in the container")
	flag.StringVar(&cgRoot, "cg_root", "/sys/fs/cgroup", "path to where cgroups are mounted")
}

func main() {
	for _, envvar := range os.Environ() {
		kv := strings.SplitN(envvar, "=", 2)
		envFlag[kv[0]] = kv[1]
	}

	// by default, set environment variables and pretend to be LXC
	// http://cgit.freedesktop.org/systemd/systemd/tree/src/shared/virt.c#n171
	envFlag["container"] = "lxc"

	flag.Parse()

	vfs, _ := lnxns.NewVfs(cgRoot)
	if iscg, err := vfs.IsCgroupFs(); !iscg {
		fmt.Printf("%s does not appear to be a cgroup filesystem: %s\n", cgRoot, err)
	}

	// create a Cgroup
	cg, _ := lnxns.NewCgroup(vfs, cgroupName)

	// add this process to the cgroup, children will inherit
	cg.AddProcess(os.Getpid())

	args := flag.Args()
	argv := make([]string, len(args)+1)
	argv[0] = programFlag
	for i := range args {
		argv[i+1] = args[i]
	}

	fmt.Printf("syscall.Exec('%s', '%s', '%s')\n", programFlag, argv, envFlag)
	err := syscall.Exec(programFlag, argv, os.Environ())
	if err != nil {
		fmt.Printf("exec failed: %s\n", err)
	}
}

// vim: ts=4 sw=4 noet tw=120 softtabstop=4
