package image

import (
	"fmt"
	"image"
	"io"

	"github.com/jonathanpittman/tiff"
)

func validateTIFF(t *tiff.TIFF) error {
	if len(t.IFDs) == 0 {
		return fmt.Errorf("tiff/image: no IFDs present in tiff to process")
	}
	if t.IFDs[0] == nil {
		return fmt.Errorf("tiff/image: IFD 0 is nil")
	}
	if t.IFDs[0].NumEntries() == 0 || len(t.IFDs[0].Fields()) == 0 {
		return fmt.Errorf("tiff/image: no entries found in IFD 0")
	}
	return nil
}

func Decode(r io.Reader) (img image.Image, err error) {
	var t *tiff.TIFF
	if t, err = tiff.ParseTIFF(tiff.NewReadAtReadSeeker(r), nil, nil); err != nil {
		return
	}
	if err = validateTIFF(t); err != nil {
		return
	}

	if handlr := findAlternates(t); handlr != nil {
		return handlr.Process(t)
	}

	// If no alternates are found...  do our own thing as best we can.
	// 1. Retrieve pixel data
	// 2. Process pixel data into either 8 or 16 bit pixels
	// 3. Make an image.Image
	err = fmt.Errorf("tiff/image: generic image decoding not yet implemented")
	return
}

func DecodeConfig(r io.Reader) (cfg image.Config, err error) {
	var t *tiff.TIFF
	if t, err = tiff.ParseTIFF(tiff.NewReadAtReadSeeker(r), nil, nil); err != nil {
		return
	}
	if err = validateTIFF(t); err != nil {
		return
	}

	if handlr := findAlternates(t); handlr != nil {
		return handlr.GetConfig(t)
	}

	err = fmt.Errorf("tiff/image: generic config decoding not yet implemented")
	return
}

func init() {
	image.RegisterFormat("tiff", tiff.TIFFBigEndian, Decode, DecodeConfig)
	image.RegisterFormat("tiff", tiff.TIFFLitEndian, Decode, DecodeConfig)
}
