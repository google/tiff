package image

import (
	"fmt"
	"image"
	"image/color"

	"github.com/jonathanpittman/tiff"
)

/* Baseline RGB

Differences from Palette Color Images

Tag 258 (BitsPerSample)
	BitsPerSample = 8,8,8
	Note: Each component is 8 bits deep in a Baseline TIFF RGB image.

Tag 262 (PhotometricInterpretation)
	PhotometricInterpretation = 2 (RGB)

There is no ColorMap.

Tag 277 (SamplesPerPixel)
	Note: The number of components per pixel. This number is 3 for RGB images, unless extra samples are present. See the ExtraSamples field for further information.

Required Fields
		<Baseline Grayscale>
		ImageWidth
		ImageLength
		BitsPerSample
		Compression
		PhotometricInterpretation
		StripOffsets
	SamplesPerPixel
		RowsPerStrip
		StripByteCounts
		XResolution
		YResolution
		ResolutionUnit

*/

type fullColorRGBDecoder struct {
	grayscaleDecoder `tiff:"ifd"`
	SamplesPerPixel  uint16 `tiff:"field,tag=277"`
}

func (rgbDec *fullColorRGBDecoder) Image() (image.Image, error) {
	if rgbDec.img == nil {
		return nil, fmt.Errorf("tiff/image: no baseline rgb image found")
	}
	return rgbDec.img, nil
}

func (rgbDec *fullColorRGBDecoder) Config() (cfg image.Config, err error) {
	cfg.Height = int(rgbDec.ImageLength)
	cfg.Width = int(rgbDec.ImageWidth)
	cfg.ColorModel = color.RGBAModel
	return
}

type FullColorRGB struct{}

func (FullColorRGB) Decoder(ifd tiff.IFD, br tiff.BReader) (dec Decoder, err error) {
	rgbDec := &fullColorRGBDecoder{grayscaleDecoder: grayscaleDecoder{bilevelDecoder: bilevelDecoder{br: br}}}
	if err = tiff.UnmarshalIFD(ifd, rgbDec); err != nil {
		return
	}
	fmt.Printf("rgb: %#v\n", rgbDec)
	return rgbDec, nil
}

func (FullColorRGB) CanHandle(ifd tiff.IFD) bool {
	return new(Grayscale).CanHandle(ifd) && ifd.HasField(277)
}
