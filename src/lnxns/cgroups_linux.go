// Copyright 2013 Albert P. Tobey. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package lnxns

import (
	"fmt"
	"io/ioutil"
	"strings"
)

type Cgroup struct {
}

func NewCgroup(name string) (cg *Cgroup) {
	return nil
}

func (cg *Cgroup) AddController() {
	return
}

func (cg *Cgroup) Controllers() (list []string) {
	return nil
}

// add a process by pid, automatically getting all threads
func (cg *Cgroup) AddProcess(tid int) {
	return
}

// add a task by tid/pid, does not recurse
func (cg *Cgroup) AddTask(tid int) {
	return
}

func (cg *Cgroup) Tasks() (list []int) {
	return
}

func (cg *Cgroup) Apply() (err error) {
	return
}

// finds where cgroups are mounted and returns the path string
func FindMountPath() (path string, err error) {
	mtab, err := ioutil.ReadFile("/proc/mounts")
	assertNil(err, "Could read /proc/mounts")

	for mount := range strings.Split(string(mtab), "\n") {
		fmt.Print(mount)
	}
	return "", err
}

// vim: ts=4 sw=4 noet tw=120 softtabstop=4
