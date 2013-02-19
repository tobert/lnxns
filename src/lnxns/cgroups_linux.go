// Copyright 2013 Albert P. Tobey. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package lnxns

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

type Cgroup struct {
	Name        string
	root        *Vfs
	isMountRoot bool // is this a mount root or unified mount?
	controllers []string
}

// create a new cgroup, the Vfs provided should already point
// to the root of cgroup mount(s).  If the cgroup already exists
// it is _not_ read automatically, you should call .Load() if you
// want to start with what the kernel has.
// root, err := FindCgroupRoot()
// root, err := NewVfs("/sys/fs/cgroup")
// root, err := NewVfs("/sys/fs/cgroup/memory")
// cg := NewCgroup(root, "tobert")
func NewCgroup(vfs *Vfs, name string) (cg *Cgroup, err error) {
	// see if we're in a controller mount or a root
	// start by looking for a "tasks" file, if it's there, vfs is either a single
	// controller mount or monolithic, either way just scan the files to see which
	// controllers are enabled and call it done. Otherwise, it's a root like /sys/fs/cgroup
	// and all path construction is a little different.

	// TODO: Write this for real ...
	taskFile := path.Join(vfs.Path(), "tasks")

	rs, err := os.Stat(taskFile)
	if err != nil {
		fmt.Printf("no such file!? %s\n", err)
		return
	}
	rs.Mode()
	return
}

// search for controllers in the root Vfs and return them
func ListAvailableControllers() (list []string, err error) {
	return
}

func (cg *Cgroup) AddController() {
	return
}

func (cg *Cgroup) Controllers() (list []string) {
	for _, vfs := range cg.controllers {
		list = append(list, string(vfs))
	}
	return
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

// load settings from an existing cgroup, returns ENOENT if the
// cgroup doesn't already exist
func (cg *Cgroup) Load() (err error) {
	return
}

// compare in-memory values to what the kernel has, returns
// a bool and a string diff (format is not guaranteed)
func (cg *Cgroup) Verify() (match bool, diff string, err error) {
	return
}

func (cg *Cgroup) Apply() (err error) {
	return
}

// finds where cgroups are mounted and returns the path string
func FindCgroupRoot() (vfs *Vfs, err error) {
	mtab, err := ioutil.ReadFile("/proc/mounts")
	assertNil(err, "Could read /proc/mounts")

	for mount := range strings.Split(string(mtab), "\n") {
		fmt.Print(mount)
	}
	return
}

// vim: ts=4 sw=4 noet tw=120 softtabstop=4
