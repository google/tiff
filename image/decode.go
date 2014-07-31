package image

import (
	"fmt"
	"image"
	"io"

	"github.com/jonathanpittman/tiff"
)

func init() {
	image.RegisterFormat("TIFF", tiff.TIFFBigEndian, DecodeTIFF, DecodeConfig)
	image.RegisterFormat("TIFF", tiff.TIFFLitEndian, DecodeTIFF, DecodeConfig)
	image.RegisterFormat("BigTIFF", tiff.BigTIFFBigEndian, DecodeBigTIFF, DecodeConfig)
	image.RegisterFormat("BigTIFF", tiff.BigTIFFLitEndian, DecodeBigTIFF, DecodeConfig)
}

func DecodeTIFF(r io.Reader) (image.Image, error) {
	/*
		rars := tiff.NewReadAtReadSeeker(r)
		t, err := tiff.ParseTIFF(rars, nil, nil)
	*/
	return nil, fmt.Errorf("Not yet implemented.")
}

func DecodeBigTIFF(r io.Reader) (image.Image, error) {
	/*
		rars := tiff.NewReadAtReadSeeker(r)
		bt, err := tiff.ParseBigTIFF(rars, nil, nil)
	*/
	return nil, fmt.Errorf("Not yet implemented.")
}

func DecodeConfig(r io.Reader) (out image.Config, err error) {
	err = fmt.Errorf("Not yet implemented.")
	return
}
