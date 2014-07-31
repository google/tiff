package tiff

import (
	"encoding/binary"
	"fmt"
)

type TIFFHeader struct {
	Order       uint16 // "MM" or "II"
	Version     uint16 // Must be 42 (0x2A)
	FirstOffset uint32 // Offset location for IFD 0
}

type TIFF struct {
	TIFFHeader
	IFDs []IFD
	R    BReader
}

func (t *TIFF) ByteOrder() binary.ByteOrder {
	return getByteOrder(t.Order)
}

func ParseTIFF(r ReadAtReadSeeker, tsp TagSpace, fts FieldTypeSet) (out *TIFF, err error) {
	if tsp == nil {
		tsp = DefaultTagSpace
	}
	if fts == nil {
		fts = DefaultFieldTypes
	}

	var th TIFFHeader

	// Get the byte order
	if err = binary.Read(r, binary.BigEndian, &th.Order); err != nil {
		return
	}
	// Check the byte order
	order := getByteOrder(th.Order)
	if order == nil {
		return nil, fmt.Errorf("tiff: invalid byte order %q", []byte{byte(th.Order >> 8), byte(th.Order)})
	}

	br := NewBReader(r, order)

	// Get the TIFF type
	if err = br.BRead(&th.Version); err != nil {
		return
	}
	// Check the type (42 for TIFF)
	if th.Version != VersionTIFF {
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

	t := &TIFF{
		TIFFHeader: th,
		R:          br,
	}

	// Locate and process IFDs
	for nextOffset := t.FirstOffset; nextOffset != 0; {
		var ifd IFD
		if ifd, err = ParseIFD(br, nextOffset, tsp, fts); err != nil {
			return
		}
		t.IFDs = append(t.IFDs, ifd)
		nextOffset = ifd.NextOffset()
	}
	return t, nil
}
