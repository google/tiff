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

package image // import "jonathanpittman.com/tiff/image"

import (
	"fmt"
	"image"
	"image/color"
	"math/big"

	"jonathanpittman.com/tiff"
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

type bilevelDecoder struct {
	ImageWidth                uint32   `tiff:"field,tag=256"`
	ImageLength               uint32   `tiff:"field,tag=257"`
	Compression               uint16   `tiff:"field,tag=259"`
	PhotometricInterpretation uint16   `tiff:"field,tag=262"`
	StripOffsets              []uint32 `tiff:"field,tag=273"`
	RowsPerStrip              uint32   `tiff:"field,tag=278"`
	StripByteCounts           []uint32 `tiff:"field,tag=279"`
	XResolution               *big.Rat `tiff:"field,tag=282"`
	YResolution               *big.Rat `tiff:"field,tag=283"`
	ResolutionUnit            uint16   `tiff:"field,tag=296"`

	br  tiff.BReader
	img image.Image
}

func (bld *bilevelDecoder) Image() (image.Image, error) {
	if bld.img == nil {
		return nil, fmt.Errorf("tiff/image: no baseline bilevel image found")
		// Do decoding here...
	}
	return bld.img, nil
}

func (bld *bilevelDecoder) Config() (cfg image.Config, err error) {
	cfg.Height = int(bld.ImageLength)
	cfg.Width = int(bld.ImageWidth)
	// Maybe we should investigate creating a B/W model
	cfg.ColorModel = color.GrayModel
	return
}

type Bilevel struct{}

func (Bilevel) Decoder(ifd tiff.IFD, br tiff.BReader) (dec Decoder, err error) {
	blDec := &bilevelDecoder{br: br}
	if err = tiff.UnmarshalIFD(ifd, blDec); err != nil {
		return
	}
	fmt.Printf("bilevel: %#v\n", blDec)
	return blDec, nil
}

func (Bilevel) CanHandle(ifd tiff.IFD) bool {
	requiredTags := []uint16{256, 257, 259, 262, 273}
	for _, tagID := range requiredTags {
		if !ifd.HasField(tagID) {
			return false
		}
	}
	return true
}
