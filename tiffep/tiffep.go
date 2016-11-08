// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tiffep // import "jonathanpittman.com/tiff/tiffep"

import "jonathanpittman.com/tiff"

var tiffEPTags = tiff.NewTagSet("TIFF/EP", 32768, 65535)

func init() {
	tiffEPTags.Register(tiff.NewTag(33421, "CFARepeatPatternDim", nil))
	tiffEPTags.Register(tiff.NewTag(33422, "CFAPattern", nil))
	tiffEPTags.Register(tiff.NewTag(34859, "SelfTimeMode", nil))
	tiffEPTags.Register(tiff.NewTag(37390, "FocalPlaneXResolution", nil))
	tiffEPTags.Register(tiff.NewTag(37391, "FocalPlaneYResolution", nil))
	tiffEPTags.Register(tiff.NewTag(37392, "FocalPlaneResolutionUnit", nil))
	tiffEPTags.Register(tiff.NewTag(37398, "TIFF/EPStandardID", nil))
	tiffEPTags.Register(tiff.NewTag(37399, "SensingMethod", nil))

	tiffEPTags.Lock()

	tiff.DefaultTagSpace.RegisterTagSet(tiffEPTags)
}
