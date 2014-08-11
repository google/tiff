// Package tiff85 provides parsing for a tiff file with a version number of 85.
package tiff85

import (
	"encoding/binary"
	"fmt"

	"github.com/jonathanpittman/tiff"
)

const (
	Version         uint16 = 0x55 // 85
	TIFF85BigEndian        = "MM\x00\x55"
	TIFF85LitEndian        = "II\x55\x00"
)

// ParseTIFF85 is practically the same as the regular tiff package's ParseTIFF
// function except for the version checking.  This functionality is kept
// separate to keep the core tiff package slim and free from non-standard bits.
func ParseTIFF85(r tiff.ReadAtReadSeeker, tsp tiff.TagSpace, ftsp tiff.FieldTypeSpace) (out *tiff.TIFF, err error) {
	if tsp == nil {
		tsp = tiff.DefaultTagSpace
	}
	if ftsp == nil {
		ftsp = tiff.DefaultFieldTypeSpace
	}

	var th tiff.Header

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
		Header: th,
		R:      br,
	}

	// Locate and process IFDs
	for nextOffset := t.FirstOffset; nextOffset != 0; {
		var ifd tiff.IFD
		if ifd, err = tiff.ParseIFD(br, nextOffset, tsp, ftsp); err != nil {
			return
		}
		t.IFDs = append(t.IFDs, ifd)
		nextOffset = ifd.NextOffset()
	}
	return t, nil
}
