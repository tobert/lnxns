// Copyright 2013 Albert P. Tobey. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package lnxns

import (
	"syscall"
)

func NsFork(more_flags int) (pid int, err error) {
	// CLONE_NEWNET unsupported for now
	// assume the caller wants an isolated process and turn on all namespacing except
	// for network, which requires setup to get networking going in the container
	// explicitly avoid CLONE_FS/CLONE_FILES/CLONE_IO and threading-related flags!
	var flags int = CLONE_NEWNS | CLONE_NEWPID | CLONE_NEWUTS | CLONE_NEWIPC | SIGCHLD | more_flags

	// see go/src/pkg/syscall/exec_unix.go
	syscall.ForkLock.Lock()

	r1, _, err1 := syscall.RawSyscall(syscall.SYS_CLONE, uintptr(flags), 0, 0)

	syscall.ForkLock.Unlock()

	if err1 != 0 {
		return 0, err1
	}

	// parent will get the pid, child will be 0
	return int(r1), nil
}

// vim: ts=4 sw=4 noet tw=120 softtabstop=4
