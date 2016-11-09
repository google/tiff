// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package exif

import "github.com/google/tiff"

const InteroperabilityIFDTagID = 40965

var (
	iopTags     = tiff.NewTagSet("Interoperability", 0, 65535)
	IOPTagSpace = tiff.NewTagSpace("Interoperability")
	iopIFDTag   = tiff.NewTag(InteroperabilityIFDTagID, "InteroperabilityIFD", nil)
)

func init() {
	tiff.PrivateTags.Register(iopIFDTag)

	// http://www.exiv2.org/tags.html
	iopTags.Register(tiff.NewTag(1, "InteroperabilityIndex", nil))
	iopTags.Register(tiff.NewTag(2, "InteroperabilityVersion", nil))
	iopTags.Register(tiff.NewTag(4096, "RelatedImageFileFormat", nil))
	iopTags.Register(tiff.NewTag(4097, "RelatedImageWidth", nil))
	iopTags.Register(tiff.NewTag(4098, "RelatedImageLength", nil))

	iopTags.Lock()

	IOPTagSpace.RegisterTagSet(iopTags)
	tiff.RegisterTagSpace(IOPTagSpace)
}
