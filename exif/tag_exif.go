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

import (
	"math/big"
	"strings"

	"jonathanpittman.com/tiff"
)

/*
http://www.awaresystems.be/imaging/tiff/tifftags/privateifd/exif.html
http://www.exiv2.org/tags.html
http://www.cipa.jp/exifprint/index_e.html
http://www.cipa.jp/std/documents/e/DC-008-2012_E.pdf
http://www.cipa.jp/std/documents/e/DC-008-Translation-2016-E.pdf
http://www.jeita.or.jp/cgi-bin/standard_e/list.cgi?cateid=1&subcateid=4
*/

const ExifIFDTagID = 34665

var (
	exifTags     = tiff.NewTagSet("Exif", 0, 65535)
	ExifTagSpace = tiff.NewTagSpace("Exif")
	exifIFDTag   = tiff.NewTag(ExifIFDTagID, "ExifIFD", nil)
)

/*
TODO: Break up these exif tags into sets based on the exif version.  They
still all likely belong in the same space though.  For an example, take
a look at the way DNG was broken up.  Tags introduced in newer versions
are added to a set named for the version.  They still all get put into
the same space, just the sets are identified separately.
*/

/*
TODO: Break up tags into categories
From TIFF:
  A. Tags relating to image data structure
  B. Tags relating to recording offset
  C. Tags relating to image data characteristics
  D. Other tags
From EXIF:
  A. Tags Relating to Version
  B. Tag Relating to Image Data Characteristics
  C. Tags Relating to Image Configuration
  D. Tags Relating to User Information
  E. Tag Relating to Related File Information
  F. Tags Relating to Date and Time
  G. Tags Relating to Picture-Taking Conditions
  G2. Tags Relating to shooting situation
  H. Other Tags
From GPS:
  A. Tags Relating to GPS
From Interoperability:
  A. Attached Information Related to Interoperability
*/

// fiRat displays rationals (fractions) in n/d notation.  It assumes the
// underlying Go type is a *big.Rat.  This can be better.
func fiRat(f tiff.Field) string {
	return f.Type().Valuer()(f.Value().Bytes(), f.Value().Order()).Interface().(*big.Rat).String()
}

// fiRatAsFloat displays a rational as a float with 2 decimal points.
func fiRatAsFloat(f tiff.Field) string {
	return f.Type().Valuer()(f.Value().Bytes(), f.Value().Order()).Interface().(*big.Rat).FloatString(2)
}

/* f/Stop formatting:
The purpose of this is to do:
  "0.95" -> "f/0.95"
  "2.80" -> "f/2.8"
  "4.00" -> "f/4"
*/
// fiFNumber displays an ƒ/stop value in the form f/n.nn, f/n.n, or f/n.
func fiFNumber(f tiff.Field) string {
	return "f/" + strings.TrimSuffix(strings.TrimRight(fiRatAsFloat(f), "0"), ".")
}

// ExposureProgram value names.  In the future, this may need to be a map if
// index values are skipped.
var exposureProgramVals = [...]string{
	"Not defined",       // 0
	"Manual",            // 1
	"Normal program",    // 2
	"Aperture priority", // 3
	"Shutter priority",  // 4
	"Creative program",  // 5 (biased toward depth of field)
	"Action program",    // 6 (biased toward fast shutter speed)
	"Portrait mode",     // 7 (for closeup photos with the background out of focus)
	"Landscape mode",    // 8 (for landscape photos with the background in focus)
}

func fiExposureProgram(f tiff.Field) string {
	prog := f.Type().Valuer()(f.Value().Bytes(), f.Value().Order()).Uint()
	if prog >= uint64(len(exposureProgramVals)) {
		return ""
	}
	return exposureProgramVals[prog]
}

func init() {
	tiff.PrivateTags.Register(exifIFDTag)

	exifTags.Register(tiff.NewTag(33434, "ExposureTime", fiRat))
	exifTags.Register(tiff.NewTag(33437, "FNumber", fiFNumber))
	exifTags.Register(tiff.NewTag(34850, "ExposureProgram", fiExposureProgram))
	exifTags.Register(tiff.NewTag(34852, "SpectralSensitivity", nil))
	exifTags.Register(tiff.NewTag(34855, "ISOSpeedRatings", nil))
	exifTags.Register(tiff.NewTag(34856, "OECF", nil))
	exifTags.Register(tiff.NewTag(34864, "SensitivityType", nil))
	exifTags.Register(tiff.NewTag(34866, "RecommendedExposureIndex", nil))
	exifTags.Register(tiff.NewTag(36864, "ExifVersion", nil))
	exifTags.Register(tiff.NewTag(36867, "DateTimeOriginal", nil))
	exifTags.Register(tiff.NewTag(36868, "DateTimeDigitized", nil))
	exifTags.Register(tiff.NewTag(37121, "ComponentsConfiguration", nil))
	exifTags.Register(tiff.NewTag(37122, "CompressedBitsPerPixel", nil))
	exifTags.Register(tiff.NewTag(37377, "ShutterSpeedValue", nil))
	exifTags.Register(tiff.NewTag(37378, "ApertureValue", nil))
	exifTags.Register(tiff.NewTag(37379, "BrightnessValue", nil))
	exifTags.Register(tiff.NewTag(37380, "ExposureBiasValue", nil))
	exifTags.Register(tiff.NewTag(37381, "MaxApertureValue", nil))
	exifTags.Register(tiff.NewTag(37382, "SubjectDistance", nil))
	exifTags.Register(tiff.NewTag(37383, "MeteringMode", nil))
	exifTags.Register(tiff.NewTag(37384, "LightSource", nil))
	exifTags.Register(tiff.NewTag(37385, "Flash", nil))
	exifTags.Register(tiff.NewTag(37386, "FocalLength", nil))
	exifTags.Register(tiff.NewTag(37396, "SubjectArea", nil))
	exifTags.Register(tiff.NewTag(37500, "MakerNote", nil))
	exifTags.Register(tiff.NewTag(37510, "UserComment", nil))
	exifTags.Register(tiff.NewTag(37520, "SubsecTime", nil))
	exifTags.Register(tiff.NewTag(37521, "SubsecTimeOriginal", nil))
	exifTags.Register(tiff.NewTag(37522, "SubsecTimeDigitized", nil))
	exifTags.Register(tiff.NewTag(40960, "FlashpixVersion", nil))
	exifTags.Register(tiff.NewTag(40961, "ColorSpace", nil))
	exifTags.Register(tiff.NewTag(40962, "PixelXDimension", nil))
	exifTags.Register(tiff.NewTag(40963, "PixelYDimension", nil))
	exifTags.Register(tiff.NewTag(40964, "RelatedSoundFile", nil))
	exifTags.Register(tiff.NewTag(41483, "FlashEnergy", nil))
	exifTags.Register(tiff.NewTag(41484, "SpatialFrequencyResponse", nil))
	exifTags.Register(tiff.NewTag(41486, "FocalPlaneXResolution", nil))
	exifTags.Register(tiff.NewTag(41487, "FocalPlaneYResolution", nil))
	exifTags.Register(tiff.NewTag(41488, "FocalPlaneResolutionUnit", nil))
	exifTags.Register(tiff.NewTag(41492, "SubjectLocation", nil))
	exifTags.Register(tiff.NewTag(41493, "ExposureIndex", nil))
	exifTags.Register(tiff.NewTag(41495, "SensingMethod", nil))
	exifTags.Register(tiff.NewTag(41728, "FileSource", nil))
	exifTags.Register(tiff.NewTag(41729, "SceneType", nil))
	exifTags.Register(tiff.NewTag(41730, "CFAPattern", nil))
	exifTags.Register(tiff.NewTag(41985, "CustomRendered", nil))
	exifTags.Register(tiff.NewTag(41986, "ExposureMode", nil))
	exifTags.Register(tiff.NewTag(41987, "WhiteBalance", nil))
	exifTags.Register(tiff.NewTag(41988, "DigitalZoomRatio", nil))
	exifTags.Register(tiff.NewTag(41989, "FocalLengthIn35mmFilm", nil))
	exifTags.Register(tiff.NewTag(41990, "SceneCaptureType", nil))
	exifTags.Register(tiff.NewTag(41991, "GainControl", nil))
	exifTags.Register(tiff.NewTag(41992, "Contrast", nil))
	exifTags.Register(tiff.NewTag(41993, "Saturation", nil))
	exifTags.Register(tiff.NewTag(41994, "Sharpness", nil))
	exifTags.Register(tiff.NewTag(41995, "DeviceSettingDescription", nil))
	exifTags.Register(tiff.NewTag(41996, "SubjectDistanceRange", nil))
	exifTags.Register(tiff.NewTag(42016, "ImageUniqueID", nil))
	exifTags.Register(tiff.NewTag(42032, "CameraOwnerName", nil))
	exifTags.Register(tiff.NewTag(42033, "BodySerialNumber", nil))
	exifTags.Register(tiff.NewTag(42034, "LensSpecification", nil))
	exifTags.Register(tiff.NewTag(42035, "LensMake", nil))
	exifTags.Register(tiff.NewTag(42036, "LensModel", nil))
	exifTags.Register(tiff.NewTag(42037, "LensSerialNumber", nil))

	// Tags that indicate the offsets to the respective IFDs.
	exifTags.Register(exifIFDTag)
	exifTags.Register(gpsIFDTag)
	exifTags.Register(iopIFDTag)

	// Not sure if this actually belongs in Exif, but it has shown up in an ExifIFD.
	exifTags.Register(tiff.NewTag(18246, "Rating", nil))

	// Prevent further registration in exif.  If tags are missing, they
	// should be added here instead of added from the outside.
	exifTags.Lock()

	tiff.DefaultTagSpace.RegisterTagSet(exifTags)

	ExifTagSpace.RegisterTagSet(tiff.BaselineTags)
	ExifTagSpace.RegisterTagSet(tiff.ExtendedTags)
	ExifTagSpace.RegisterTagSet(exifTags)

	tiff.RegisterTagSpace(ExifTagSpace)
}
