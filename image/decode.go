package image

import (
	"fmt"
	"image"
	"io"
)

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
