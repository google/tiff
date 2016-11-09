// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package modi provides tiff extensions for working with Microsoft Office
// Document Imaging based tiff files.
package modi

import "github.com/google/tiff"

var modiTags = tiff.NewTagSet("MODI", 32768, 65535)

func init() {
	modiTags.Register(tiff.NewTag(37679, "MODIText", nil))
	modiTags.Register(tiff.NewTag(37680, "MODIOLEPropertySetStorage", nil))
	modiTags.Register(tiff.NewTag(37681, "MODIPositioning", nil))

	modiTags.Lock()

	tiff.DefaultTagSpace.RegisterTagSet(modiTags)
}
