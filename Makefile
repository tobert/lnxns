# Copyright 2013 Albert P. Tobey. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

GO := GOPATH=$(shell pwd) go

all: test binaries

fmt:
	$(shell cd src/lnxns && go fmt)
	$(shell cd nschroot  && go fmt)
	$(shell cd cgroup    && go fmt)
	$(shell cd contain   && go fmt)

test:
	$(GO) test ./src/lnxns

binaries:
	$(GO) build -o nschroot/nschroot nschroot/main.go
	$(GO) build -o cgroup/cgroup cgroup/main.go
	$(GO) build -o contain/contain contain/main.go

clean:
	rm -f nschroot/nschroot cgroup/cgroup contain/contain

# vim: ts=4 sw=4 noet tw=120 softtabstop=4
