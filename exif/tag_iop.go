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

package exif // import "jonathanpittman.com/tiff/exif"

import "jonathanpittman.com/tiff"

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
