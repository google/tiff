// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package image // import "jonathanpittman.com/tiff/image"

import (
	"fmt"
	"image"
	"image/color"
	"math/big"

	"jonathanpittman.com/tiff"
)

type Baseline struct {
	NewSubfileType            *uint32  `tiff:"field,tag=254"`
	SubfileType               *uint16  `tiff:"field,tag=255"`
	ImageWidth                *uint32  `tiff:"field,tag=256"`
	ImageLength               *uint32  `tiff:"field,tag=257"`
	BitsPerSample             []uint16 `tiff:"field,tag=258"`
	Compression               *uint16  `tiff:"field,tag=259"`
	PhotometricInterpretation *uint16  `tiff:"field,tag=262"`
	Threshholding             *uint16  `tiff:"field,tag=263"`
	CellWidth                 *uint16  `tiff:"field,tag=264"`
	CellLength                *uint16  `tiff:"field,tag=265"`
	FillOrder                 *uint16  `tiff:"field,tag=266"`
	ImageDescription          *string  `tiff:"field,tag=270"`
	Make                      *string  `tiff:"field,tag=271"`
	Model                     *string  `tiff:"field,tag=272"`
	StripOffsets              []uint32 `tiff:"field,tag=273"`
	Orientation               *uint16  `tiff:"field,tag=274"`
	SamplesPerPixel           *uint16  `tiff:"field,tag=277"`
	RowsPerStrip              *uint32  `tiff:"field,tag=278"`
	StripByteCounts           []uint32 `tiff:"field,tag=279"`
	MinSampleValue            *uint16  `tiff:"field,tag=280"`
	MaxSampleValue            *uint16  `tiff:"field,tag=281"`
	XResolution               *big.Rat `tiff:"field,tag=282"`
	YResolution               *big.Rat `tiff:"field,tag=283"`
	PlanarConfiguration       *uint16  `tiff:"field,tag=284"`
	FreeOffsets               []uint32 `tiff:"field,tag=288"`
	FreeByteCounts            []uint32 `tiff:"field,tag=289"`
	GrayResponseUnit          *uint16  `tiff:"field,tag=290"`
	GrayResponseCurve         []uint16 `tiff:"field,tag=291"`
	ResolutionUnit            *uint16  `tiff:"field,tag=296"`
	Software                  *string  `tiff:"field,tag=305"`
	DateTime                  *string  `tiff:"field,tag=306"`
	Artist                    *string  `tiff:"field,tag=315"`
	HostComputer              *string  `tiff:"field,tag=316"`
	ColorMap                  []uint16 `tiff:"field,tag=320"`
	ExtraSamples              []byte   `tiff:"field,tag=338"`
	Copyright                 *string  `tiff:"field,tag=33432"`
}

// Required fields for baseline
type BaselineDecoder struct {
	//We only really care about these fields.
	ImageWidth                *uint32  `tiff:"field,tag=256"`
	ImageLength               *uint32  `tiff:"field,tag=257"`
	BitsPerSample             []uint16 `tiff:"field,tag=258"`
	Compression               *uint16  `tiff:"field,tag=259"`
	PhotometricInterpretation *uint16  `tiff:"field,tag=262"`
	StripOffsets              []uint32 `tiff:"field,tag=273"`
	SamplesPerPixel           *uint16  `tiff:"field,tag=277"`
	RowsPerStrip              *uint32  `tiff:"field,tag=278"`
	StripByteCounts           []uint32 `tiff:"field,tag=279"`
	XResolution               *big.Rat `tiff:"field,tag=282"`
	YResolution               *big.Rat `tiff:"field,tag=283"`
	ResolutionUnit            *uint16  `tiff:"field,tag=296"`
	ColorMap                  []uint16 `tiff:"field,tag=320"`

	bl  *Baseline `tiff:"ifd,idx=0"`
	br  tiff.BReader
	img image.Image
}

func (bld *BaselineDecoder) Image() (image.Image, error) {
	if bld.img == nil {
		return nil, fmt.Errorf("tiff/image: no baseline bilevel image found")
		// Do decoding here...
	}
	return bld.img, nil
}

func (bld *BaselineDecoder) Config() (cfg image.Config, err error) {
	if bld.ImageLength == nil {
		err = fmt.Errorf("tiff/image: missing value for ImageLength")
		return
	}
	cfg.Height = int(*bld.ImageLength)
	if bld.ImageWidth == nil {
		err = fmt.Errorf("tiff/image: missing value for ImageWidth")
		return
	}
	cfg.Width = int(*bld.ImageWidth)
	switch {
	case len(bld.ColorMap) > 0:
		// palette color
		// TODO: Create a palette color ColorModel from the ColorMap field.
	case bld.SamplesPerPixel != nil && *bld.SamplesPerPixel > 0:
		switch *bld.SamplesPerPixel {
		case 3:
			cfg.ColorModel = color.RGBAModel
			for _, bps := range bld.BitsPerSample {
				if bps > 8 {
					cfg.ColorModel = color.RGBA64Model
					break
				}
			}
		case 1:
			cfg.ColorModel = color.GrayModel
			for _, bps := range bld.BitsPerSample {
				if bps > 8 {
					cfg.ColorModel = color.Gray16Model
					break
				}
			}
		default:
			err = fmt.Errorf("tiff/image: unsupported SamplesPerPixel value: %d", bld.SamplesPerPixel)
		}
	}

	return
}

type BaselineHandler struct{}

func (BaselineHandler) Decoder(ifd tiff.IFD, br tiff.BReader) (dec Decoder, err error) {
	return nil, err
	/*
		switch {
		case ifd.HasField(320): // ColorMap (PaletteColor)
			return new(PaletteColor).Decoder(ifd, br)
		case ifd.HasField(277): // SamplesPerPixel (RGB)
			sppField := ifd.GetField(277)
			spp := sppField.Value().Order().Uint16(sppField.Value().Bytes())
			switch spp {
			case 1:
				if ifd.HasField(258) {
					bpsField := ifd.GetField(258)
					bps := bpsField.Value().Order().Uint16(bpsField.Value().Bytes())
					if bps > 1 {
						return new(Grayscale).Decoder(ifd, br)
					}
				}
				return new(Bilevel).Decoder(ifd, br)
			case 3:
				return new(FullColorRGB).Decoder(ifd, br)
			}
			return nil, fmt.Errorf("tiff/image: unsupported SamplesPerPixel value: %d", spp)
		case ifd.HasField(258): // BitsPerSample (GrayScale)
			return new(Grayscale).Decoder(ifd, br)
		}
		// Assume bilevel at the very least
		return new(Bilevel).Decoder(ifd, br)
	*/
}

func (BaselineHandler) CanHandle(ifd tiff.IFD) bool {
	return new(Bilevel).CanHandle(ifd)
}
