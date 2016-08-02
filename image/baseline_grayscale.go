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

	"jonathanpittman.com/tiff"
)

/* Baseline grayscaleDecoder

Differences from Bilevel Images
	Change:
		Tag 259 (Compression)
			Note: Must be 1 or 32773 only

	Add:
		Tag 258 (BitsPerSample)
			4 = 16 shades of gray
			8 = 256 shades of gray
			Note: The number of bits per component.
			Note: Allowable values for Baseline TIFF grayscale
				images are 4 and 8.

Required Fields
		<Baseline Bilevel>
		ImageWidth
		ImageLength
	BitsPerSample
		Compression
		PhotometricInterpretation
		StripOffsets
		RowsPerStrip
		StripByteCounts
		XResolution
		YResolution
		ResolutionUnit
*/

type grayscaleDecoder struct {
	bilevelDecoder `tiff:"ifd"`
	BitsPerSample  []uint16 `tiff:"field,tag=258"`
}

func (gsd *grayscaleDecoder) Image() (image.Image, error) {
	if gsd.img == nil {
		return nil, fmt.Errorf("tiff/image: no baseline grayscale image found")
	}
	return gsd.img, nil
}

func (gsd *grayscaleDecoder) Config() (cfg image.Config, err error) {
	cfg.Height = int(gsd.ImageLength)
	cfg.Width = int(gsd.ImageWidth)
	cfg.ColorModel = color.GrayModel
	return
}

type Grayscale struct{}

func (Grayscale) Decoder(ifd tiff.IFD, br tiff.BReader) (dec Decoder, err error) {
	gsDec := &grayscaleDecoder{bilevelDecoder: bilevelDecoder{br: br}}
	if err = tiff.UnmarshalIFD(ifd, gsDec); err != nil {
		return
	}
	fmt.Printf("grayscale: %#v\n", gsDec)
	return gsDec, nil
}

func (Grayscale) CanHandle(ifd tiff.IFD) bool {
	return new(Bilevel).CanHandle(ifd) && ifd.HasField(258)
}
