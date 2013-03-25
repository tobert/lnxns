# lnxns - Linux namespaces in Go

## Warning

This is brand-new software. I've tested it minimally before pushing to Github. Do not
expect signatures and layout to be consistent until this note disappears.

## Requirements

Linux >= 2.6.24 with:

    CONFIG_NAMESPACES=y
    CONFIG_UTS_NS=y
    CONFIG_IPC_NS=y
    CONFIG_PID_NS=y
    CONFIG_NET_NS=y (eventually)

Root or CAP_SYS_ADMIN privileges. Using setcap on a binary may not be safe on a multi-user
system since input checking isn't very thorough.

## Build

    make
    make test
    make clean
    make binaries

## Example

If busybox is installed, this should work

    sudo ./nschroot /bin /busybox ls /

    mkdir -p /tmp/root
    cp -a /bin/busybox /tmp/root
    touch /tmp/root/foobar
    go build -o nschroot nschroot.go && sudo ./nschroot /tmp/root /busybox ls

## History

* 2013-03-25: 'nschroot' and 'cgroup' are working
* 2013-02-19: nschroot seems to work fine as root. Cgroups aren't there yet, but I should have a workable API soon.

## Author

Al Tobey <tobert@gmail.com> @AlTobey

## License

Copyright 2013 Albert P Tobey.  All rights reserved.
Use of this source code is governed by a BSD-style
license that can be found in the LICENSE file.
