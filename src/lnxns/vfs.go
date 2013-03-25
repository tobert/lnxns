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
type Vfs struct {
	Device     string
	Mountpoint string
	Filesystem string
	Options    []string
}

// create a new Vfs handle, checks that it exists, does not verify
// if it's a valid fs beyond making sure it's a directory, since it's
// useful to test against temp dirs or even use fake vfs (fuse)
func NewVfs(mpath string) (*Vfs, error) {
	// if it's a mounted path, return with all options set via /proc/mounts
	mtab := Mounts()
	if _, ok := mtab[mpath]; ok {
		return mtab[mpath], nil
	}

	// otherwise, the user will have to set other fields as needed
	var vfs Vfs = Vfs{Mountpoint: mpath}

	rs, err := os.Stat(mpath)
	if err != nil {
		return nil, err
	}

	if !rs.Mode().IsDir() {
		return nil, errors.New("not a directory")
	}

	return &vfs, nil
}

// return the VFS path as a string
func (vfs *Vfs) Path() string {
	return vfs.Mountpoint
}

// returns a Vfs set up for working with /proc, assuming it's good to go
func ProcFs() *Vfs {
	proc := Vfs{
		Device:     "proc",
		Mountpoint: "/proc",
		Filesystem: "proc",
		Options:    []string{"rw"},
	}

	return &proc
}

func SysFs() *Vfs {
	sys, err := NewVfs("/sys")
	if err != nil {
		panic(fmt.Sprintf("/sys unavailable: %s", err))
	}
	return sys
}

// parse /proc/mounts and return a map of mountpoint: *Vfs
func Mounts() map[string]*Vfs {
	var ret = make(map[string]*Vfs)

	mtab, err := ProcFs().GetMapList("mounts", 1)
	assertNil(err, "Could read /proc/mounts")

	for mp, opts := range mtab {
		v := Vfs{
			Device:     opts[0],
			Mountpoint: opts[1],
			Filesystem: opts[2],
			Options:    strings.Split(opts[3], ","),
		}
		ret[mp] = &v
	}

	return ret
}

// check if the Vfs is pointing at an instance of proc
// A proc Vfs must point at the root of the proc mountpoint, e.g. "/proc".
func (vfs *Vfs) IsProcFs() (bool, error) {
	for _, required := range procfsRequires {
		st, err := os.Stat(path.Join(vfs.Mountpoint, required))
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
		st, err := os.Stat(path.Join(vfs.Mountpoint, required))
		if err != nil || !st.Mode().IsDir() {
			return false, err
		}
	}

	return true, nil
}

// check if the Vfs is pointing at some kind of cgroup fs
func (vfs *Vfs) IsCgroupFs() (bool, error) {
	mtab := Mounts()

	if _, ok := mtab[vfs.Mountpoint]; ok {
		switch mtab[vfs.Mountpoint].Filesystem {
		// systemd style mounts each controller under a tmpfs in /sys/fs/cgroup
		case "tmpfs":
			// if cpuset is all there and has a tasks file, call it good enough
			if _, err := os.Stat(path.Join(vfs.Mountpoint, "cpuset", "tasks")); err == nil {
				return true, nil
			}
		// but some people like to mount it monolithic, e.g. mount -t cgroup none /cgroups
		case "cgroup":
			if _, err := os.Stat(path.Join(vfs.Mountpoint, "tasks")); err == nil {
				return true, nil
			} else {
				panic("Invalid cgroup filesystem! All cgroup mountpoints must have a 'tasks' file.")
			}
		default:
			return false, errors.New("not a supported filesystem")
		}

	}

	return false, errors.New("Does not appear to be a mountpoint.")
}

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
		pt   string = path.Join(vfs.Mountpoint, name)
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
		pt   string = path.Join(vfs.Mountpoint, name)
	)

	if file, err = os.Create(pt); err != nil {
		return
	}
	defer file.Close()

	_, err = fmt.Fprintln(file, value)
	return
}

// vim: ts=4 sw=4 noet tw=120 softtabstop=4
