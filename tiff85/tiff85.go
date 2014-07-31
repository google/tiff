package tiff85

import (
	"encoding/binary"
	"fmt"

	"github.com/jonathanpittman/tiff"
)

const (
	Version       uint16 = 0x55
	TIFFBigEndian        = "MM\x00\x55"
	TIFFLitEndian        = "II\x55\x00"
)

func ParseTIFF85(r tiff.ReadAtReadSeeker, tsp tiff.TagSpace, fts tiff.FieldTypeSet) (out *tiff.TIFF, err error) {
	if tsp == nil {
		tsp = tiff.DefaultTagSpace
	}
	if fts == nil {
		fts = tiff.DefaultFieldTypes
	}

	var th tiff.TIFFHeader

	// Get the byte order
	if err = binary.Read(r, binary.BigEndian, &th.Order); err != nil {
		return
	}
	// Check the byte order
	order := tiff.GetByteOrder(th.Order)
	if order == nil {
		return nil, fmt.Errorf("tiff: invalid byte order %q", []byte{byte(th.Order >> 8), byte(th.Order)})
	}

	br := tiff.NewBReader(r, order)

	// Get the TIFF type
	if err = br.BRead(&th.Version); err != nil {
		return
	}
	// Check the type (85 for this version of TIFF)
	if th.Version != Version {
		return nil, fmt.Errorf("tiff: unsupported version %d", th.Version)
	}

	// Get the offset to the first IFD
	if err = br.BRead(&th.FirstOffset); err != nil {
		return
	}
	// Check the offset to the first IFD (ensure it is past the end of the header)
	if th.FirstOffset < 8 {
		return nil, fmt.Errorf("tiff: invalid offset to first IFD, %d < 8", th.FirstOffset)
	}

	t := &tiff.TIFF{
		TIFFHeader: th,
		R:          br,
	}

	// Locate and process IFDs
	for nextOffset := t.FirstOffset; nextOffset != 0; {
		var ifd tiff.IFD
		if ifd, err = tiff.ParseIFD(br, nextOffset, tsp, fts); err != nil {
			return
		}
		t.IFDs = append(t.IFDs, ifd)
		nextOffset = ifd.NextOffset()
	}
	return t, nil
}
