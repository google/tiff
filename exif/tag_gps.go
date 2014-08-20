package exif

import (
	"math/big"

	"github.com/jonathanpittman/tiff"
)

// TODO: Pass in the slice of valid FieldType for each tag.

const GPSIFDTagID = 34853

var (
	gpsTags     = tiff.NewTagSet("GPS", 0, 65535)
	GPSTagSpace = tiff.NewTagSpace("GPS")
	gpsIFDTag   = tiff.NewTag(GPSIFDTagID, "GPS IFD", nil)
)

func init() {
	tiff.PrivateTags.Register(gpsIFDTag)

	// http://www.awaresystems.be/imaging/tiff/tifftags/privateifd/gps.html
	// http://www.exiv2.org/tags.html
	gpsTags.Register(tiff.NewTag(0, "GPSVersionID", nil))
	gpsTags.Register(tiff.NewTag(1, "GPSLatitudeRef", nil))
	gpsTags.Register(tiff.NewTag(2, "GPSLatitude", nil))
	gpsTags.Register(tiff.NewTag(3, "GPSLongitudeRef", nil))
	gpsTags.Register(tiff.NewTag(4, "GPSLongitude", nil))
	gpsTags.Register(tiff.NewTag(5, "GPSAltitudeRef", nil))
	gpsTags.Register(tiff.NewTag(6, "GPSAltitude", nil))
	gpsTags.Register(tiff.NewTag(7, "GPSTimeStamp", nil))
	gpsTags.Register(tiff.NewTag(8, "GPSSatellites", nil))
	gpsTags.Register(tiff.NewTag(9, "GPSStatus", nil))
	gpsTags.Register(tiff.NewTag(10, "GPSMeasureMode", nil))
	gpsTags.Register(tiff.NewTag(11, "GPSDOP", nil))
	gpsTags.Register(tiff.NewTag(12, "GPSSpeedRef", nil))
	gpsTags.Register(tiff.NewTag(13, "GPSSpeed", nil))
	gpsTags.Register(tiff.NewTag(14, "GPSTrackRef", nil))
	gpsTags.Register(tiff.NewTag(15, "GPSTrack", nil))
	gpsTags.Register(tiff.NewTag(16, "GPSImgDirectionRef", nil))
	gpsTags.Register(tiff.NewTag(17, "GPSImgDirection", nil))
	gpsTags.Register(tiff.NewTag(18, "GPSMapDatum", nil))
	gpsTags.Register(tiff.NewTag(19, "GPSDestLatitudeRef", nil))
	gpsTags.Register(tiff.NewTag(20, "GPSDestLatitude", nil))
	gpsTags.Register(tiff.NewTag(21, "GPSDestLongitudeRef", nil))
	gpsTags.Register(tiff.NewTag(22, "GPSDestLongitude", nil))
	gpsTags.Register(tiff.NewTag(23, "GPSDestBearingRef", nil))
	gpsTags.Register(tiff.NewTag(24, "GPSDestBearing", nil))
	gpsTags.Register(tiff.NewTag(25, "GPSDestDistanceRef", nil))
	gpsTags.Register(tiff.NewTag(26, "GPSDestDistance", nil))
	gpsTags.Register(tiff.NewTag(27, "GPSProcessingMethod", nil))
	gpsTags.Register(tiff.NewTag(28, "GPSAreaInformation", nil))
	gpsTags.Register(tiff.NewTag(29, "GPSDateStamp", nil))
	gpsTags.Register(tiff.NewTag(30, "GPSDifferential", nil))

	gpsTags.Lock()

	GPSTagSpace.RegisterTagSet(gpsTags)
	tiff.RegisterTagSpace(GPSTagSpace)
}

type gpsIFD struct {
	VersionID        [4]byte     `tifftag:"id=0"`
	LatitudeRef      string      `tifftag:"id=1,type=2"`
	Latitude         [3]*big.Rat `tifftag:"id=2"`
	LongitudeRef     string      `tifftag:"id=3"`
	Longitude        [3]*big.Rat `tifftag:"id=4"`
	AltitudeRef      byte        `tifftag:"id=5"`
	Altitude         *big.Rat    `tifftag:"id=6"`
	TimeStamp        [3]*big.Rat `tifftag:"id=7"`
	Satellites       string      `tifftag:"id=8"`
	Status           string      `tifftag:"id=9"`
	MeasureMode      string      `tifftag:"id=10"`
	DOP              *big.Rat    `tifftag:"id=11"`
	SpeedRef         string      `tifftag:"id=12"`
	Speed            *big.Rat    `tifftag:"id=13"`
	TrackRef         string      `tifftag:"id=14"`
	Track            *big.Rat    `tifftag:"id=15"`
	ImgDirectionRef  string      `tifftag:"id=16"`
	ImgDirection     *big.Rat    `tifftag:"id=17"`
	MapDatum         string      `tifftag:"id=18"`
	DestLatitudeRef  string      `tifftag:"id=19"`
	DestLatitude     *big.Rat    `tifftag:"id=20"`
	DestLongitudeRef string      `tifftag:"id=21"`
	DestLongitude    *big.Rat    `tifftag:"id=22"`
	DestBearingRef   string      `tifftag:"id=23"`
	DestBearing      *big.Rat    `tifftag:"id=24"`
	DestDistanceRef  string      `tifftag:"id=25"`
	DestDistance     *big.Rat    `tifftag:"id=26"`
	ProcessingMethod []byte      `tifftag:"id=27"`
	AreaInformation  []byte      `tifftag:"id=28"`
	DateStamp        string      `tifftag:"id=29"`
	Differential     uint16      `tifftag:"id=30"`
}
