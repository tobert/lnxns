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

type Vfs string

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
	parser := func(parts []string) {
		var key string = parts[keyIndex]
		values[key] = parts
	}

	vfs.slurp(name, parser)

	return
}

func (vfs *Vfs) SetString(name string, value string) (err error) {
	return vfs.write(name, value)
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
			for _, part := range strings.Split(line, " ") {
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
