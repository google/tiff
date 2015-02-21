package image

import (
	"fmt"
	"image"
	"image/color"

	"github.com/jonathanpittman/tiff"
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
