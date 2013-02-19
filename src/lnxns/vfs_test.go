// Copyright 2013 Albert P. Tobey. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package lnxns_test

import (
	"../../src/lnxns"
	"os"
	"path"
	"testing"
)

// type Mountpoint string
// func NewVfs(mpath string) (vr *Mountpoint, err error)
// func (vr *Vfs) GetString(name string) (value string, err error)
// func (vr *Vfs) GetInt(name string) (value int, err error)
// func (vr *Vfs) GetIntList(name string) (values []int, err error)
// func (vr *Vfs) GetMapList(name string, keyIndex int) (values map[string][]string, err error)

func TestNewVfs(t *testing.T) {
	vr, err := lnxns.NewVfs("/0abc1def2ghi3jkl4mno5pqr6stu7vwx8yz9")
	if err == nil {
		t.Fatalf("NewVfs returned a nil error where a real error was expected.")
	}
	if vr != nil {
		t.Fatalf("NewVfs returned a Vfs when it shouldn't have.")
	}

	tmpDir := os.TempDir()
	// TODO: find a better tmpdir function
	tmpPath := path.Join(tmpDir, "test-lnxns")
	os.Mkdir(tmpPath, 0755)
	defer os.RemoveAll(tmpPath)

	vr, err = lnxns.NewVfs(tmpPath)
	if err != nil {
		t.Fatalf("NewVfs %q: %s", tmpPath, err)
	}

	foo, err := vr.GetString("foo")
	if foo != "" {
		t.Fatalf("GetString returned a value for a non-existent file. Expected '', Got '%s'", foo)
	}
	if err == nil {
		t.Fatalf("GetString returned a nil error where a real error was expected.")
	}

	err = vr.SetString("int1", "9999")
	if err != nil {
		t.Fatalf("SetString returned an error! '%s'", err)
	}

	baz, err := vr.GetString("int1")
	if err != nil {
		t.Fatalf("GetString returned an error! '%s'", err)
	}
	if baz != "9999" {
		t.Fatalf("GetString failed, Got: '%s'", baz)
	}

	bar, err := vr.GetInt("int1")
	if err != nil {
		t.Fatalf("GetInt returned an error! '%s'", err)
	}
	if bar != 9999 {
		t.Fatalf("GetInt failed, Got: '%d'", bar)
	}
}

// vim: ts=4 sw=4 noet tw=120 softtabstop=4
