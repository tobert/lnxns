// Copyright 2013 Albert P. Tobey. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package lnxns_test

import (
	"../../src/lnxns"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"testing"
)

func TestCgroup(t *testing.T) {
	tmpPath, _ := ioutil.TempDir(os.TempDir(), "test-lnxns-cgroups")
	os.Mkdir(tmpPath, 0755)
	defer os.RemoveAll(tmpPath)

	os.Mkdir(path.Join(tmpPath, "memory"), 0755)
	os.Mkdir(path.Join(tmpPath, "blkio"), 0755)

	vr, _ := lnxns.NewVfs(tmpPath)
	vr.SetString("memory/memory.swappiness", "0")
	vr.SetString("memory/tasks", "123\n456\n789")
}

func TestFindCgroups(t *testing.T) {
	vfs := lnxns.FindCgroupVfs()
	fmt.Printf("VFS: %s\n", vfs)
}

// vim: ts=4 sw=4 noet tw=120 softtabstop=4
