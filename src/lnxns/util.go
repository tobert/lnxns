// Copyright 2013 Albert P. Tobey. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package lnxns

import (
	"fmt"
	"log"
)

func assertNil(err error, message string) {
	if err != nil {
		log.Fatal(fmt.Sprintf("%s: %s\n", message, err))
	}
}

// vim: ts=4 sw=4 noet tw=120 softtabstop=4
