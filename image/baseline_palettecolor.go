package image // import "jonathanpittman.com/tiff/image"

import (
	"fmt"
	"image"
	"image/color"

	"jonathanpittman.com/tiff"
)

/* Baseline Palette-Color

Differences from Grayscale Images
	Tag 262 (PhotometricInterpretation)
		PhotometricInterpretation = 3 (Palette Color)

	Tag 320 (ColorMap)
		Note: 3 * (2**BitsPerSample)
		Note: This field defines a Red-Green-Blue color map (often
			called a lookup table) for palette color images. In a
			palette-color image, a pixel value is used to index into
			an RGB-lookup table. For example, a palette-color pixel
			having a value of 0 would be displayed according to the
			0th Red, Green, Blue triplet.
			In a TIFF ColorMap, all the Red values come first,
			followed by the Green values, then the Blue values. In
			the ColorMap, black is represented by 0,0,0 and white is
			represented by 65535, 65535, 65535.

Required Fields
		<Baseline Grayscale>
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
	ColorMap
*/

type paletteColorDecoder struct {
	grayscaleDecoder `tiff:"ifd"`
	ColorMap         []uint16 `tiff:"field,tag=320"`
}

func (pcd *paletteColorDecoder) Image() (image.Image, error) {
	if pcd.img == nil {
		return nil, fmt.Errorf("tiff/image: no baseline palette color image found")
	}
	return pcd.img, nil
}

func (pcd *paletteColorDecoder) Config() (cfg image.Config, err error) {
	cfg.Height = int(pcd.ImageLength)
	cfg.Width = int(pcd.ImageWidth)
	cfg.ColorModel = color.Palette{} // TODO: Take this from the ColorMap field
	return
}

type PaletteColor struct{}

func (PaletteColor) Decoder(ifd tiff.IFD, br tiff.BReader) (dec Decoder, err error) {
	pcDec := &paletteColorDecoder{grayscaleDecoder: grayscaleDecoder{bilevelDecoder: bilevelDecoder{br: br}}}
	if err = tiff.UnmarshalIFD(ifd, pcDec); err != nil {
		return
	}
	return pcDec, nil
}

func (PaletteColor) CanHandle(ifd tiff.IFD) bool {
	return new(Grayscale).CanHandle(ifd) && ifd.HasField(320)
}
