// Copyright 2013 Albert P. Tobey. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package lnxns

// from /usr/include/linux/sched.h
const (
	CLONE_FS      = 0x00000200 /* set if fs info shared between processes */
	CLONE_FILES   = 0x00000400 /* set if open files shared between processes */
	CLONE_NEWNS   = 0x00020000 /* New namespace group? */
	CLONE_NEWUTS  = 0x04000000 /* New utsname group? */
	CLONE_NEWIPC  = 0x08000000 /* New ipcs */
	CLONE_NEWUSER = 0x10000000 /* New user namespace */
	CLONE_NEWPID  = 0x20000000 /* New pid namespace */
	CLONE_NEWNET  = 0x40000000 /* New network namespace */
	CLONE_IO      = 0x80000000 /* Clone io context */
	CLONE_VFORK   = 0x00004000 /* set if the parent wants the child to wake it up on mm_release */
	SIGCHLD       = 0x14       /* Should set SIGCHLD for fork()-like behavior on Linux */
)

// vim: ts=4 sw=4 noet tw=120 softtabstop=4
