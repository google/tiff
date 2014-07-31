package bigtiff

import (
	"encoding/binary"
	"fmt"

	"github.com/jonathanpittman/tiff"
)

type BigTIFFHeader struct {
	Order       uint16 // "MM" or "II"
	Version     uint16 // Must be 43 (0x2B)
	OffsetSize  uint16 // Size in bytes used for offset values
	Constant    uint16 // Must be 0
	FirstOffset uint64 // Offset location for IFD 0
}

type BigTIFF struct {
	BigTIFFHeader
	IFDs []IFD8
	R    tiff.BReader
}

func (bt *BigTIFF) ByteOrder() binary.ByteOrder {
	return tiff.GetByteOrder(bt.Order)
}

func ParseBigTIFF(r tiff.ReadAtReadSeeker, tsg tiff.TagSpace, fts tiff.FieldTypeSet) (out *BigTIFF, err error) {
	if tsg == nil {
		tsg = tiff.DefaultTagSpace
	}
	if fts == nil {
		fts = tiff.DefaultFieldTypes
	}

	var bth BigTIFFHeader

	// Get the byte order
	if err = binary.Read(r, binary.BigEndian, &bth.Order); err != nil {
		return
	}
	// Check the byte order
	order := tiff.GetByteOrder(bth.Order)
	if order == nil {
		return nil, fmt.Errorf("tiff: invalid byte order %q", []byte{byte(bth.Order >> 8), byte(bth.Order)})
	}

	br := tiff.NewBReader(r, order)

	// Get the TIFF type
	if err = br.BRead(&bth.Version); err != nil {
		return
	}
	// Check the type (43 for BigTIFF)
	if bth.Version != Version {
		return nil, fmt.Errorf("tiff: unsupported version %d", bth.Version)
	}

	// Get the offset size
	if err = br.BRead(&bth.OffsetSize); err != nil {
		return
	}
	// Check the offset size (For now, only support an offset size of 8 for uint64.)
	if bth.OffsetSize != 8 {
		return nil, fmt.Errorf("tiff: invalid offset size of %d", bth.OffsetSize)
	}

	// Get the constant
	if err = br.BRead(&bth.Constant); err != nil {
		return
	}
	// Check the constant
	if bth.Constant != 0 {
		return nil, fmt.Errorf("tiff: invalid header constant, %d != 0", bth.Constant)
	}

	// Get the offset to the first IFD
	if err = br.BRead(&bth.FirstOffset); err != nil {
		return
	}
	// Check the offset to the first IFD (ensure it is past the end of the header)
	if bth.FirstOffset < 16 {
		return nil, fmt.Errorf("tiff: invalid offset to first IFD, %d < 16", bth.FirstOffset)
	}

	bt := &BigTIFF{
		BigTIFFHeader: bth,
		R:             br,
	}

	// Locate and process IFDs
	for nextOffset := bt.FirstOffset; nextOffset != 0; {
		var ifd IFD8
		if ifd, err = ParseIFD8(br, nextOffset, tsg, fts); err != nil {
			return
		}
		bt.IFDs = append(bt.IFDs, ifd)
		nextOffset = ifd.NextOffset()
	}
	return bt, nil
}
