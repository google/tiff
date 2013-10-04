package tiff

import (
	"image"
	"io"
)

// RegisterFormat(name, magic string, decode func(io.Reader) (Image, error), decodeConfig func(io.Reader) (Config, error))

func init() {
	image.RegisterFormat("TIFF", hdrTIFFBigEndian, decodeTIFF, decodeConfig)
	image.RegisterFormat("TIFF", hdrTIFFLitEndian, decodeTIFF, decodeConfig)
	image.RegisterFormat("BigTIFF", hdrBigTIFFBigEndian, decodeBigTIFF, decodeConfig)
	image.RegisterFormat("BigTIFF", hdrBigTIFFLitEndian, decodeBigTIFF, decodeConfig)
}

func decodeTIFF(r io.Reader) (image.Image, error) {
	rars := NewReadAtReadSeeker(r)
	t, err := ParseTIFF(rars)
	_ = t.IFDs[0].ImageData
	return nil, err
}

func decodeBigTIFF(r io.Reader) (image.Image, error) {
	rars := NewReadAtReadSeeker(r)
	bt, err := ParseBigTIFF(rars)
	_ = bt.IFDs[0].ImageData
	return nil, err
}

func decodeConfig(r io.Reader) (out image.Config, err error) {
	return
}
