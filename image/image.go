package image

import (
	"fmt"
	"image"
	"io"

	"github.com/jonathanpittman/tiff"

	"github.com/jonathanpittman/tiff/bigtiff"
	"github.com/jonathanpittman/tiff/tiff85"
)

type Decoder interface {
	Image() (image.Image, error)
	Config() (image.Config, error)
}

type TIFFHandler interface {
	Decoder(tiff.TIFF) (Decoder, error)
	CanHandle(tiff.TIFF) bool
}

type IFDHandler interface {
	Decoder(tiff.IFD, tiff.BReader) (Decoder, error)
	CanHandle(tiff.IFD) bool
}

func validateTIFF(t tiff.TIFF) error {
	if len(t.IFDs()) == 0 {
		return fmt.Errorf("tiff/image: no IFDs present in tiff to process")
	}
	if t.IFDs()[0] == nil {
		return fmt.Errorf("tiff/image: IFD 0 is nil")
	}
	if t.IFDs()[0].NumEntries() == 0 || len(t.IFDs()[0].Fields()) == 0 {
		return fmt.Errorf("tiff/image: no entries found in IFD 0")
	}
	return nil
}

func getDecoder(t tiff.TIFF) (dec Decoder, err error) {
	if err = validateTIFF(t); err != nil {
		return
	}

	// Look for alternates that can handle the whole tiff.
	if handlr := findAlternateTIFFHandler(t); handlr != nil {
		return handlr.Decoder(t)
	}

	// Look for alternates that can handle specific IFDs.
	for _, ifd := range t.IFDs() {
		if handlr := findAlternateIFDHandler(ifd); handlr != nil {
			return handlr.Decoder(ifd, t.R())
		}
	}

	// If no alternates are found...  do our own thing as best we can, which
	// means baseline support only.
	ifd0 := t.IFDs()[0]
	if !new(BaselineHandler).CanHandle(ifd0) {
		return nil, fmt.Errorf("tiff/image: no handlers available for this tiff")
	}
	return new(BaselineHandler).Decoder(ifd0, t.R())
}

func Decode(r io.Reader) (img image.Image, err error) {
	var dec Decoder
	var t tiff.TIFF
	if t, err = tiff.Parse(tiff.NewReadAtReadSeeker(r), nil, nil); err != nil {
		return
	}
	if dec, err = getDecoder(t); err != nil {
		return
	}
	return dec.Image()
}

func DecodeConfig(r io.Reader) (cfg image.Config, err error) {
	var dec Decoder
	var t tiff.TIFF
	if t, err = tiff.Parse(tiff.NewReadAtReadSeeker(r), nil, nil); err != nil {
		return
	}
	if dec, err = getDecoder(t); err != nil {
		return
	}
	return dec.Config()
}

func init() {
	image.RegisterFormat("tiff", tiff.MagicBigEndian, Decode, DecodeConfig)
	image.RegisterFormat("tiff", tiff.MagicLitEndian, Decode, DecodeConfig)
	image.RegisterFormat("bigtiff", bigtiff.MagicBigEndian, Decode, DecodeConfig)
	image.RegisterFormat("bigtiff", bigtiff.MagicLitEndian, Decode, DecodeConfig)
	image.RegisterFormat("tiff85", tiff85.MagicBigEndian, Decode, DecodeConfig)
	image.RegisterFormat("tiff85", tiff85.MagicLitEndian, Decode, DecodeConfig)
}
