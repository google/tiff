package image

import (
	"fmt"
	"image"
	"math/big"

	"github.com/jonathanpittman/tiff"
)

/* Baseline Bilevel

Color
	Tag 262 (PhotometricInterpretation)
		0 = WhiteIsZero (normal when Compression=2)
		1 = BlackIsZero (If Compression=2, the image should display and print reversed)

Compression
	Tag 259 (Compression)
		1 = uncompressed (when encoding/writing, pack data into bytes tightly)
		2 = CCITT Group 3 1-Dimensional Modified Huffman run length encoding (CCITT_Grp3_1-D_MH_RLE)
		32773 = PackBits
		Note: Data compression applies only to raster image data. All other TIFF fields are unaffected.
		Note: Baseline TIFF readers must handle all three compression schemes.

Rows and Columns
	Tag 257 (ImageLength)
		Note: Number of rows

	Tag 256 (ImageWidth)
		Note: Number of Columns

Physical Dimensions
	Tag 296 (ResolutionUnit)
		1 = No absolute unit of measurement
		2 = inch (default)
		3 = centimeter

	Tag 282 (XResolution)
		Type: Rational
		Note: number of pixels per ResolutionUnit in the ImageWidth

	Tag 283 (YResolution)
		Type: Rational
		Note: number of pixels per ResolutionUnit in the ImageLength

Location of the Data
	Tag 278 (RowsPerStrip)
		Note: The number of rows in each strip (except possibly the last strip.)

	Tag 273 (StripOffsets)
		Note: For each strip, the byte offset of that strip.

	Tag 279 (StripByteCounts)
		Note: For each strip, the number of bytes in that strip after any compression.


Required Fields
	ImageWidth
	ImageLength
	Compression
	PhotometricInterpretation
	StripOffsets
	RowsPerStrip
	StripByteCounts
	XResolution
	YResolution
	ResolutionUnit
*/

type Bilevel struct {
	ImageWidth                uint32     `tifftag:"id=256"`
	ImageLength               uint32     `tifftag:"id=257"`
	Compression               uint16     `tifftag:"id=259"`
	PhotometricInterpretation uint16     `tifftag:"id=262"`
	StripOffsets              []uint32   `tifftag:"id=273"`
	RowsPerStrip              uint32     `tifftag:"id=278"`
	StripByteCounts           []uint32   `tifftag:"id=279"`
	XResolution               *big.Rat   `tifftag:"id=282"`
	YResolution               *big.Rat   `tifftag:"id=283"`
	ResolutionUnit            uint16     `tifftag:"id=296"`
	Rest                      []tiff.Tag // Any left over tags from an IFD.
}

func (bl *Bilevel) Process(tbr tiff.BReader) (image.Image, error) {
	return nil, fmt.Errorf("tiff: bilevel handling not yet implemented")
}
