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
	"strconv"
)

type Cgroup struct {
	Name        string
	vfs			*Vfs
}

// create a new cgroup, the Vfs provided should already point to the root of
// cgroup mount(s).  If the cgroup already exists it is not read automatically,
// you should call .Load() if you want to start with what the kernel has.
// Only systemd-style per-controller mounts are supported at this time.
// cg := NewCgroup(FindCgroupVfs(), "tobert")
func NewCgroup(v *Vfs, name string) (*Cgroup, error) {
	cg := Cgroup{
		Name: name,
		vfs: v,
	}

	// if the tasks file exists, this is either a monolithic mount or a single controller
	taskFile := path.Join(v.Path(), "tasks")
	_, err := os.Stat(taskFile)
	if err == nil {
		panic(fmt.Sprintf("Found a tasks file in %s: Monolithic and single controller mounts are not supported. Try /sys/fs/cgroup.", taskFile))
	}

	for _, ctl := range ListControllers() {
		err = os.Mkdir(path.Join(v.Path(), ctl, name), 0755)
		// ignore EEXIST, it's fine and common
		if os.IsExist(err) {
			continue
		} else if err != nil {
			fmt.Printf("Could not create directory: %s\n", err)
		}
	}

	return &cg, nil
}

// Returns a list of available cgroups in the running host kernel. Reads /proc/cgroups.
// e.g. [net_cls blkio devices cpuset cpuacct memory freezer cpu]
func ListControllers() (list []string) {
	_, err := os.Stat("/proc/cgroups")
	if err != nil {
		panic("Could not stat /proc/cgroups. Your kernel does not seem to support cgroups.")
	}

	rows, err := ProcFs().GetMapList("cgroups", 0)
	if err != nil {
		panic("BUG: Could not parse /proc/cgroups.")
	}

	for key, _ := range rows {
		if strings.HasPrefix(key, "#") {
			continue
		} else {
			list = append(list, key)
		}
	}

	return list
}

// moves tasks back to the global group and deletes the directory
func (cg *Cgroup) Destroy() (error) {
	for _, ctl := range ListControllers() {
		tasks, err := cg.vfs.GetIntList(path.Join(ctl, cg.Name, "tasks"))

		// move tasks back to the root controller
		for _, task := range tasks {
			cg.vfs.SetString(path.Join(ctl, "tasks"), strconv.Itoa(task))
			// check errors?
		}

		// remove the control group
		err = os.Remove(path.Join(cg.vfs.Path(), ctl, cg.Name))
		if err != nil {
			fmt.Printf("Could not remove directory: %s\n", err)
		}
	}

	return nil
}

// return the full path to a control file. Even though the controller prefix is redundant,
// it's required since not all files in a control directory have a prefix.
// e.g.
// cg := lnxns.NewCgroup(FindCgroupVfs(), "junk")
// cg.ctlPath("memory", "memory.swappiness") == "/sys/fs/cgroup/memory/junk/memory.swappiness"
func (cg *Cgroup) ctlPath(controller string, file string) (p string) {
	p = path.Join(cg.vfs.Path(), controller, cg.Name, file)

	// TODO: remove this debug output & stat someday
	_, err := os.Stat(p)
	if err != nil {
		fmt.Printf("File '%s' does not exist: %s", p, err)
	}

	return p
}

// add a process by pid, automatically getting all threads
func (cg *Cgroup) AddProcess(pid int) {
	for _, name := range ListControllers() {
		taskFile := path.Join(name, "tasks")
		cg.vfs.SetString(taskFile, strconv.Itoa(pid))

		pid_tasks, _ := ioutil.ReadDir(path.Join("/proc", strconv.Itoa(pid), "task"))
		for _, fi := range pid_tasks {
			cg.vfs.SetString(taskFile, fi.Name())
		}
	}
}

// add a task by tid/pid, does not recurse
func (cg *Cgroup) AddTask(tid int) {
	for _, ctl := range ListControllers() {
		taskFile := cg.ctlPath(ctl, "tasks")
		cg.vfs.SetString(taskFile, strconv.Itoa(tid))
	}
}

// finds where cgroups are mounted and returns the path string
// /sys/fs/cgroup is tried first, then search /proc/mounts
func FindCgroupVfs() *Vfs {
	v, err := NewVfs("/sys/fs/cgroup")
	if err == nil && v.Filesystem == "tmpfs" {
		if iscg, _ := v.IsCgroupFs(); iscg {
			return v
		}
	}

	mtab := Mounts()
	for mp, vfs := range mtab {
		if vfs.Filesystem == "cgroup" {
			parent := path.Base(mp)
			v, err = NewVfs(parent)
			if err == nil {
				// TODO: test this somewhere ... my machines are all systemd
				// all my Ubuntu machines are also modified to mount cgroups in the systemd style
				return v
			}
		}
	}

	return nil
}

// vim: ts=4 sw=4 noet tw=120 softtabstop=4
