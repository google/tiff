package tiff

import (
	"encoding/binary"
)

// These constants represents the first 4 bytes of the file for each kind of
// TIFF along with each byte ordering.  This is mostly useful for registration
// with the "image" package from the Go standard library.
const (
	TIFFBigEndian    = "MM\x00\x2A"
	TIFFLitEndian    = "II\x2A\x00"
	BigTIFFBigEndian = "MM\x00\x2B"
	BigTIFFLitEndian = "II\x2B\x00"
)

// These constants represent the byte order options present at the beginning of
// a TIFF file.
const (
	BigEndian uint16 = 0x4D4D // "MM" or 19789
	LitEndian uint16 = 0x4949 // "II" or 18761
)

func getByteOrder(bo uint16) binary.ByteOrder {
	switch bo {
	case BigEndian:
		return binary.BigEndian
	case LitEndian:
		return binary.LittleEndian
	}
	return nil
}

// These constants represent the TIFF file type identifiers.  At present,
// there are values for a TIFF and a BigTIFF.
const (
	VersionTIFF    uint16 = 0x2A
	VersionBigTIFF uint16 = 0x2B
)
