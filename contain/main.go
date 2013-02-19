// Copyright 2013 Albert P. Tobey. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	//"../src/lnxns"
	"flag"
	"fmt"
	"strings"
)

type containerOpts struct {
	Cgroups []string
	Program []string
	Env     map[string]string
}

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
	return fmt.Sprint(env)
}

func (env envMap) Set(value string) error {
	kv := strings.SplitN(value, "=", 2)
	env[kv[0]] = kv[1]
	fmt.Printf("KV: %s\n", kv)
	return nil
}

var cgroupsFlag cgroupList
var programFlag string
var argumentsFlag string
var envFlag envMap = make(envMap)

func init() {
	flag.Var(&cgroupsFlag, "cgroup", "comma-separated list of cgroups")
	flag.Var(&envFlag, "env", "key=value environment variables")
	flag.StringVar(&programFlag, "program", "", "the program to run in the container")
}

func main() {
	flag.Parse()
}

// vim: ts=4 sw=4 noet tw=120 softtabstop=4
