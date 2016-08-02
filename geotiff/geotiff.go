/*
Copyright 2016 Google Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package geotiff // import "jonathanpittman.com/tiff/geotiff"

import "jonathanpittman.com/tiff"

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
