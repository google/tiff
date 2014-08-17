package bigtiff

import (
	"encoding/binary"
	"fmt"

	"github.com/jonathanpittman/tiff"
)

type ErrUnsuppTIFFVersion struct {
	Version uint16
}

func (eutv ErrUnsuppTIFFVersion) Error() string {
	return fmt.Sprintf("bigtiff: unsupported version %d", eutv.Version)
}

type Header struct {
	Order       uint16 // "MM" or "II"
	Version     uint16 // Must be 43 (0x2B)
	OffsetSize  uint16 // Size in bytes used for offset values
	Constant    uint16 // Must be 0
	FirstOffset uint64 // Offset location for IFD 0
}

type BigTIFF struct {
	Header
	IFDs []IFD
	R    tiff.BReader
}

func (bt *BigTIFF) ByteOrder() binary.ByteOrder {
	return tiff.GetByteOrder(bt.Order)
}

func ParseBigTIFF(r tiff.ReadAtReadSeeker, tsp tiff.TagSpace, ftsp tiff.FieldTypeSpace) (out *BigTIFF, err error) {
	if tsp == nil {
		tsp = tiff.DefaultTagSpace
	}
	if ftsp == nil {
		ftsp = tiff.DefaultFieldTypeSpace
	}

	var hdr Header

	// Get the byte order
	if err = binary.Read(r, binary.BigEndian, &hdr.Order); err != nil {
		return
	}
	// Check the byte order
	order := tiff.GetByteOrder(hdr.Order)
	if order == nil {
		return nil, fmt.Errorf("tiff: invalid byte order %q", []byte{byte(hdr.Order >> 8), byte(hdr.Order)})
	}

	br := tiff.NewBReader(r, order)

	// Get the TIFF type
	if err = br.BRead(&hdr.Version); err != nil {
		return
	}
	// Check the type (43 for BigTIFF)
	if hdr.Version != Version {
		return nil, ErrUnsuppTIFFVersion{hdr.Version}
	}

	// Get the offset size
	if err = br.BRead(&hdr.OffsetSize); err != nil {
		return
	}
	// Check the offset size (For now, only support an offset size of 8 for uint64.)
	if hdr.OffsetSize != 8 {
		return nil, fmt.Errorf("tiff: invalid offset size of %d", hdr.OffsetSize)
	}

	// Get the constant
	if err = br.BRead(&hdr.Constant); err != nil {
		return
	}
	// Check the constant
	if hdr.Constant != 0 {
		return nil, fmt.Errorf("tiff: invalid header constant, %d != 0", hdr.Constant)
	}

	// Get the offset to the first IFD
	if err = br.BRead(&hdr.FirstOffset); err != nil {
		return
	}
	// Check the offset to the first IFD (ensure it is past the end of the header)
	if hdr.FirstOffset < 16 {
		return nil, fmt.Errorf("tiff: invalid offset to first IFD, %d < 16", hdr.FirstOffset)
	}

	bt := &BigTIFF{
		Header: hdr,
		R:      br,
	}

	// Locate and process IFDs
	for nextOffset := bt.FirstOffset; nextOffset != 0; {
		var ifd IFD
		if ifd, err = ParseIFD(br, nextOffset, tsp, ftsp); err != nil {
			return
		}
		bt.IFDs = append(bt.IFDs, ifd)
		nextOffset = ifd.NextOffset()
	}
	return bt, nil
}
