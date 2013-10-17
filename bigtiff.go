package tiff

import (
	"encoding/binary"
	"fmt"
)

type BigTIFF struct {
	byteOrder   uint16 // "MM" or "II"
	Type        uint16 // Must be 43 (0x2B)
	OffsetSize  uint16 // Size in bytes used for offset values
	Constant    uint16 // Must be 0
	FirstOffset uint64 // Offset location for IFD 0
	IFDs        []*IFD8
}

func (bt *BigTIFF) ByteOrder() binary.ByteOrder {
	return getByteOrder(bt.byteOrder)
}

func ParseBigTIFF(r ReadAtReadSeeker) (out *BigTIFF, err error) {
	bt := new(BigTIFF)
	br := &bReader{
		order: binary.BigEndian,
		r:     r,
	}
	// Check the byte order
	if err = br.Read(&bt.byteOrder); err != nil {
		return
	}
	br.order = bt.ByteOrder()
	if br.order == nil {
		return nil, fmt.Errorf("tiff: invalid byte order %q", []byte{byte(bt.byteOrder >> 8), byte(bt.byteOrder)})
	}
	// Check the type (43 for BigTIFF)
	if err = br.Read(&bt.Type); err != nil {
		return
	}
	if bt.Type != TypeBigTIFF {
		return nil, fmt.Errorf("tiff: invalid type %d", bt.Type)
	}
	// Get the offset size
	if err = br.Read(&bt.OffsetSize); err != nil {
		return
	}
	// For now, only support an offset size of 8 for uint64.
	if bt.OffsetSize != 8 {
		return nil, fmt.Errorf("tiff: invalid offset size of %d", bt.OffsetSize)
	}
	// Get the constant
	if err = br.Read(&bt.Constant); err != nil {
		return
	}
	if bt.Constant != 0 {
		return nil, fmt.Errorf("tiff: invalid header constant, %d != 0", bt.Constant)
	}
	// Get the offset to the first IFD
	if err = br.Read(&bt.FirstOffset); err != nil {
		return
	}
	if bt.FirstOffset < 16 {
		return nil, fmt.Errorf("tiff: invalid offset to first IFD, %d < 16", bt.FirstOffset)
	}
	// Locate and process IFDs
	for nextOffset := bt.FirstOffset; nextOffset != 0; {
		var ifd *IFD8
		if ifd, err = parseIFD8(br, nextOffset); err != nil {
			return
		}
		if err = ifd.processImageData(br); err != nil {
			return
		}

		bt.IFDs = append(bt.IFDs, ifd)
		nextOffset = ifd.NextOffset
	}

	return bt, nil
}
