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

// Package modi provides tiff extensions for working with Microsoft Office
// Document Imaging based tiff files.
package modi // import "jonathanpittman.com/tiff/modi"

import "jonathanpittman.com/tiff"

var modiTags = tiff.NewTagSet("MODI", 32768, 65535)

func init() {
	modiTags.Register(tiff.NewTag(37679, "MODIText", nil))
	modiTags.Register(tiff.NewTag(37680, "MODIOLEPropertySetStorage", nil))
	modiTags.Register(tiff.NewTag(37681, "MODIPositioning", nil))

	modiTags.Lock()

	tiff.DefaultTagSpace.RegisterTagSet(modiTags)
}
