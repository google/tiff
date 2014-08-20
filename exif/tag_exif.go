package exif

import "github.com/jonathanpittman/tiff"

const ExifIFDTagID = 34665

var (
	exifTags     = tiff.NewTagSet("Exif", 0, 65535)
	ExifTagSpace = tiff.NewTagSpace("Exif")
	exifIFDTag   = tiff.NewTag(ExifIFDTagID, "ExifIFD", nil)
)

// TODO: Pass in the slice of valid FieldType for each tag.
func init() {
	tiff.PrivateTags.Register(exifIFDTag)

	// http://www.awaresystems.be/imaging/tiff/tifftags/privateifd/exif.html
	exifTags.Register(tiff.NewTag(33434, "ExposureTime", nil))
	exifTags.Register(tiff.NewTag(33437, "FNumber", nil))
	exifTags.Register(tiff.NewTag(34850, "ExposureProgram", nil))
	exifTags.Register(tiff.NewTag(34852, "SpectralSensitivity", nil))
	exifTags.Register(tiff.NewTag(34855, "ISOSpeedRatings", nil))
	exifTags.Register(tiff.NewTag(34856, "OECF", nil))
	exifTags.Register(tiff.NewTag(34864, "SensitivityType", nil))
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
