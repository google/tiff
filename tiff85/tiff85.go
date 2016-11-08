// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package tiff85 provides parsing for a tiff file with a version number of 85.
package tiff85 // import "jonathanpittman.com/tiff/tiff85"

import "jonathanpittman.com/tiff"

const (
	MagicBigEndian        = "MM\x00\x55"
	MagicLitEndian        = "II\x55\x00"
	Version        uint16 = 0x55 // 85
	VersionName    string = "TIFF85"
)

func init() {
	tiff.RegisterVersion(Version, tiff.ParseTIFF)
}
