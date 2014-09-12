package exif

import (
	"math/big"

	"github.com/jonathanpittman/tiff"
)

const GPSIFDTagID = 34853

var (
	gpsTags     = tiff.NewTagSet("GPS", 0, 65535)
	GPSTagSpace = tiff.NewTagSpace("GPS")
	gpsIFDTag   = tiff.NewTag(GPSIFDTagID, "GPSIFD", nil)
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
	VersionID        [4]byte     `tiff:"field,tag=0"`
	LatitudeRef      string      `tiff:"field,tag=1"`
	Latitude         [3]*big.Rat `tiff:"field,tag=2"`
	LongitudeRef     string      `tiff:"field,tag=3"`
	Longitude        [3]*big.Rat `tiff:"field,tag=4"`
	AltitudeRef      byte        `tiff:"field,tag=5"`
	Altitude         *big.Rat    `tiff:"field,tag=6"`
	TimeStamp        [3]*big.Rat `tiff:"field,tag=7"`
	Satellites       string      `tiff:"field,tag=8"`
	Status           string      `tiff:"field,tag=9"`
	MeasureMode      string      `tiff:"field,tag=10"`
	DOP              *big.Rat    `tiff:"field,tag=11"`
	SpeedRef         string      `tiff:"field,tag=12"`
	Speed            *big.Rat    `tiff:"field,tag=13"`
	TrackRef         string      `tiff:"field,tag=14"`
	Track            *big.Rat    `tiff:"field,tag=15"`
	ImgDirectionRef  string      `tiff:"field,tag=16"`
	ImgDirection     *big.Rat    `tiff:"field,tag=17"`
	MapDatum         string      `tiff:"field,tag=18"`
	DestLatitudeRef  string      `tiff:"field,tag=19"`
	DestLatitude     *big.Rat    `tiff:"field,tag=20"`
	DestLongitudeRef string      `tiff:"field,tag=21"`
	DestLongitude    *big.Rat    `tiff:"field,tag=22"`
	DestBearingRef   string      `tiff:"field,tag=23"`
	DestBearing      *big.Rat    `tiff:"field,tag=24"`
	DestDistanceRef  string      `tiff:"field,tag=25"`
	DestDistance     *big.Rat    `tiff:"field,tag=26"`
	ProcessingMethod []byte      `tiff:"field,tag=27"`
	AreaInformation  []byte      `tiff:"field,tag=28"`
	DateStamp        string      `tiff:"field,tag=29"`
	Differential     uint16      `tiff:"field,tag=30"`
}
