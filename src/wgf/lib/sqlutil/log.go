// Copyright 2014 The Wgf Authors. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package sqlutil

var Logger func(log string)

func log(body string) {
	if nil!=Logger {
		Logger(body)
	}
}
