package tiff

import (
	"encoding/binary"
	"fmt"
)

type Header struct {
	Order       uint16 // "MM" or "II"
	Version     uint16 // Must be 42 (0x2A)
	FirstOffset uint32 // Offset location for IFD 0
}

type TIFF struct {
	Header
	IFDs []IFD
	R    BReader
}

func (t *TIFF) ByteOrder() binary.ByteOrder {
	return GetByteOrder(t.Order)
}

func ParseTIFF(r ReadAtReadSeeker, tsp TagSpace, ftsp FieldTypeSpace) (out *TIFF, err error) {
	if tsp == nil {
		tsp = DefaultTagSpace
	}
	if ftsp == nil {
		ftsp = DefaultFieldTypeSpace
	}

	var hdr Header

	// Get the byte order
	if err = binary.Read(r, binary.BigEndian, &hdr.Order); err != nil {
		return
	}
	// Check the byte order
	order := GetByteOrder(hdr.Order)
	if order == nil {
		return nil, fmt.Errorf("tiff: invalid byte order %q", []byte{byte(hdr.Order >> 8), byte(hdr.Order)})
	}

	br := NewBReader(r, order)

	// Get the TIFF type
	if err = br.BRead(&hdr.Version); err != nil {
		return
	}
	// Check the type (42 for TIFF)
	if hdr.Version != Version {
		return nil, fmt.Errorf("tiff: unsupported version %d", hdr.Version)
	}

	// Get the offset to the first IFD
	if err = br.BRead(&hdr.FirstOffset); err != nil {
		return
	}
	// Check the offset to the first IFD (ensure it is past the end of the header)
	if hdr.FirstOffset < 8 {
		return nil, fmt.Errorf("tiff: invalid offset to first IFD, %d < 8", hdr.FirstOffset)
	}

	t := &TIFF{
		Header: hdr,
		R:      br,
	}

	// Locate and process IFDs
	for nextOffset := t.FirstOffset; nextOffset != 0; {
		var ifd IFD
		if ifd, err = ParseIFD(br, nextOffset, tsp, ftsp); err != nil {
			return
		}
		t.IFDs = append(t.IFDs, ifd)
		nextOffset = ifd.NextOffset()
	}
	return t, nil
}
