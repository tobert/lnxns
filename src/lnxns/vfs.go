// Copyright 2013 Albert P. Tobey. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package lnxns

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"strconv"
	"strings"
)

var sysfsRequires = []string{"bus", "class", "dev", "devices", "fs"}
var procfsRequires = []string{
	"cmdline", "cpuinfo", "filesystems", "loadavg", "meminfo",
	"self", "stat", "uptime", "version", "vmstat",
}

// Vfs is for interacting with virtual filesystems like /proc, /sys,
// configfs, and cgroups.
type Vfs string

// create a new Vfs handle, checks that it exists, does not verify
func NewVfs(mpath string) (*Vfs, error) {
	var vfs Vfs = Vfs(mpath)

	rs, err := os.Stat(mpath)
	if err != nil {
		return nil, err
	}

	if !rs.Mode().IsDir() {
		return nil, errors.New("not a directory")
	}

	return &vfs, nil
}

func ProcFs() *Vfs {
	proc, err := NewVfs("/proc")
	if err != nil {
		panic(fmt.Sprintf("/proc unavailable: %s", err))
	}
	return proc
}

// return the VFS path as a string
func (vfs *Vfs) Path() string {
	return string(*vfs)
}

// check if the Vfs is pointing at an instance of proc
// A proc Vfs must point at the root of the proc mountpoint, e.g. "/proc".
func (vfs *Vfs) IsProcFs() (isProc bool, err error) {
	for _, required := range procfsRequires {
		st, err := os.Stat(path.Join(vfs.Path(), required))
		if err != nil || !st.Mode().IsDir() {
			return false, err
		}
	}

	return true, nil
}

// check if the Vfs is pointing at an instance of sysfs
// A sysfs Vfs must point at the root of the sysfs mountpoint, e.g. "/sys".
func (vfs *Vfs) IsSysFs() (bool, error) {
	for _, required := range sysfsRequires {
		st, err := os.Stat(path.Join(vfs.Path(), required))
		if err != nil || !st.Mode().IsDir() {
			return false, err
		}
	}

	return true, nil
}

// check if the Vfs is pointing at some kind of cgroup fs
//func (vfs *Vfs) IsCgroupFs() (isProc bool, err error) {
//}

// read a parameter as a string, this will work for any of the files
// e.g. vfs.GetString("sys/net/ipv4/tcp_congestion_control") = "cubic"
func (vfs *Vfs) GetString(name string) (value string, err error) {
	parser := func(parts []string) {
		value = parts[0]
	}

	err = vfs.slurp(name, parser)
	return
}

// read a cgroup parameter that is expected to be a single integer in a file
// e.g. vfs.GetInt("sys/vm/dirty_ratio") = 10
func (vfs *Vfs) GetInt(name string) (value int, err error) {
	parser := func(parts []string) {
		var num int
		num, err = strconv.Atoi(parts[0])
		if err == nil {
			value = num
		}
	}

	err = vfs.slurp(name, parser)
	return
}

// read a list of integers
// e.g. cgvfs.GetIntList("memory/tasks") = [ 1, 200, ... ]
func (vfs *Vfs) GetIntList(name string) (values []int, err error) {
	parser := func(parts []string) {
		var num int
		num, err = strconv.Atoi(parts[0])
		if err == nil {
			values = append(values, num)
		}
	}

	vfs.slurp(name, parser)
	return
}

// get a map[string][]string where the keyIndex item on a line is the key and every other
// item, split by whitespace, is put in an array of values. The keyIndex is not deleted
// from the list.
func (vfs *Vfs) GetMapList(name string, keyIndex int) (values map[string][]string, err error) {
	values = make(map[string][]string)

	parser := func(parts []string) {
		var key string = parts[keyIndex]
		values[key] = parts
	}

	vfs.slurp(name, parser)

	return
}

// write a string
func (vfs *Vfs) SetString(name string, value string) (err error) {
	return vfs.write(name, value)
}

// list directories in the root of the Vfs
func (vfs *Vfs) Dirs() ([]string, error) {
	// TODO: THIS IS A STUB
	fmt.Printf("Stub! Vfs.Dirs()\n")
	return []string{"memory", "cpu"}, nil
}

func (vfs *Vfs) Files() ([]string, error) {
	// TODO: THIS IS A STUB
	fmt.Printf("Stub! Vfs.Files()\n")
	return []string{"cgroup.procs", "cpuset.mem_hardwall", "cpuset.memory_spread_page", "cpuset.sched_relax_domain_level", "tasks"}, nil
}

// read a file line-by-line calling the provided function for each line
func (vfs *Vfs) slurp(name string, cb func([]string)) (err error) {
	var (
		file *os.File
		pt   string = path.Join(string(*vfs), name)
	)

	if file, err = os.Open(pt); err != nil {
		return
	}
	defer file.Close()

	reader := bufio.NewReader(file)

	// hmm this could be fmt.Scan*?
	for {
		var line string
		line, err = reader.ReadString('\n')
		if line != "" {
			var parts []string
			for _, part := range strings.Fields(line) {
				parts = append(parts, strings.TrimSpace(part))
			}
			cb(parts)
		}
		if err == io.EOF {
			err = nil
			break
		}
		if err != nil {
			break
		}
	}
	return
}

func (vfs *Vfs) write(name string, value string) (err error) {
	var (
		file *os.File
		pt   string = path.Join(string(*vfs), name)
	)

	if file, err = os.Create(pt); err != nil {
		return
	}
	defer file.Close()

	_, err = fmt.Fprintln(file, value)
	return
}

// vim: ts=4 sw=4 noet tw=120 softtabstop=4
