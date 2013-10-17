package tiff

import (
	"encoding/binary"
	"fmt"
)

type TIFF struct {
	byteOrder   uint16 // "MM" or "II"
	Type        uint16 // Must be 42 (0x2A)
	FirstOffset uint32 // Offset location for IFD 0
	IFDs        []*IFD
}

func (t *TIFF) ByteOrder() binary.ByteOrder {
	return getByteOrder(t.byteOrder)
}

func ParseTIFF(r ReadAtReadSeeker) (out *TIFF, err error) {
	t := new(TIFF)
	br := &bReader{
		order: binary.BigEndian,
		r:     r,
	}
	// Check the byte order
	if err = br.Read(&t.byteOrder); err != nil {
		return
	}
	br.order = t.ByteOrder()
	if br.order == nil {
		return nil, fmt.Errorf("tiff: invalid byte order %q", []byte{byte(t.byteOrder >> 8), byte(t.byteOrder)})
	}
	// Check the type (42 for TIFF)
	if err = br.Read(&t.Type); err != nil {
		return
	}
	if t.Type != TypeTIFF {
		return nil, fmt.Errorf("tiff: invalid type %d", t.Type)
	}
	// Get the offset to the first IFD
	if err = br.Read(&t.FirstOffset); err != nil {
		return
	}
	if t.FirstOffset < 8 {
		return nil, fmt.Errorf("tiff: invalid offset to first IFD, %d < 8", t.FirstOffset)
	}
	// Locate and process IFDs
	for nextOffset := t.FirstOffset; nextOffset != 0; {
		var ifd *IFD
		if ifd, err = parseIFD(br, nextOffset); err != nil {
			return
		}
		if err = ifd.processImageData(br); err != nil {
			return
		}
		t.IFDs = append(t.IFDs, ifd)
		nextOffset = ifd.NextOffset
	}
	return t, nil
}
