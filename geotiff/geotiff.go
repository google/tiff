// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package geotiff

import "github.com/google/tiff"

var geotiffTags = tiff.NewTagSet("GeoTIFF", 32768, 65535)

func init() {
	geotiffTags.Register(tiff.NewTag(33550, "ModelPixelScaleTag", nil))
	geotiffTags.Register(tiff.NewTag(34264, "ModelTransformationTag", nil))
	geotiffTags.Register(tiff.NewTag(33922, "ModelTiepointTag", nil))
	geotiffTags.Register(tiff.NewTag(34735, "GeoKeyDirectoryTag", nil))
	geotiffTags.Register(tiff.NewTag(34736, "GeoDoubleParamsTag", nil))
	geotiffTags.Register(tiff.NewTag(34737, "GeoAsciiParamsTag", nil))
	geotiffTags.Register(tiff.NewTag(33920, "IntergraphIrasBMatrixTag", nil))

	geotiffTags.Lock()

	tiff.DefaultTagSpace.RegisterTagSet(geotiffTags)
}
