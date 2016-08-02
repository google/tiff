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
