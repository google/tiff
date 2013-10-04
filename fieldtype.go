package tiff

import (
	"encoding/binary"
)

/* Field type definitions
From [TIFF6]:
	1  = BYTE	8-bit unsigned integer.
	2  = ASCII	8-bit byte that contains a 7-bit ASCII code; the last byte must be NUL (binary zero).
	3  = SHORT	16-bit (2-byte) unsigned integer.
	4  = LONG	32-bit (4-byte) unsigned integer.
	5  = RATIONAL	Two LONGs: the first represents the numerator of a fraction; the second, the denominator.
	6  = SBYTE	An 8-bit signed (twos-complement) integer.
	7  = UNDEFINED	An 8-bit byte that may contain anything, depending on the definition of the field.
	8  = SSHORT	A 16-bit (2-byte) signed (twos-complement) integer.
	9  = SLONG	A 32-bit (4-byte) signed (twos-complement) integer.
	10 = SRATIONAL	Two SLONG’s: the first represents the numerator of a fraction, the second the denominator.
	11 = FLOAT	Single precision (4-byte) IEEE format.
	12 = DOUBLE	Double precision (8-byte) IEEE format.
From [BIGTIFFDESIGN]:
	13 = IFD	?? 32-bit unsigned integer offset value ??
	14 = UNICODE	??
	15 = COMPLEX	??
	16 = LONG8	64-bit unsigned integer.
	17 = SLONG8	64-bit signed integer.
	18 = IFD8	64-bit unsigned integer offset value
*/

// A FieldType represents all of the necessary pieces of information one needs to
// know about a field type including a function that knows how to represent that
// type of data in an often human readable string.  Other string representation
// formats could be implemented (json, xml, etc).  Field types themselves have no
// actual stored value inside a TIFF.  They are here to help an implementer or
// user understand their format.
type FieldType interface {
	Id() uint16
	Name() string
	Size() uint64
	Signed() bool
	Repr() func([]byte, binary.ByteOrder) string
}

type fieldType struct {
	id     uint16
	name   string
	size   uint64
	signed bool
	repr   func([]byte, binary.ByteOrder) string
}

func (ft *fieldType) Id() uint16 {
	return ft.id
}

func (ft *fieldType) Name() string {
	return ft.name
}

func (ft *fieldType) Size() uint64 {
	return ft.size
}

func (ft *fieldType) Signed() bool {
	return ft.signed
}

func (ft *fieldType) Repr() func([]byte, binary.ByteOrder) string {
	return ft.repr
}

/*
Default set of Field types
  1-12:  Field types 1 - 12 are described in [TIFF6].
  13-15: Field Type IDs 13 - 15 are mentioned in [BIGTIFFDESIGN], but only to
         explain why values 13 - 15 were skipped when identifying new Field Types
         for BigTIFF. These are meant to be used with regular TIFF, but were
         apparently not properly documented prior to the BigTIFF design
         discussion.
  16-18: Field Type IDs 16 - 18 were added for use with BigTIFF in [BIGTIFFDESIGN].
*/
var (
	fTByte      = &fieldType{id: 1, name: "BYTE", size: 1, signed: false, repr: nil}
	fTASCII     = &fieldType{id: 2, name: "ASCII", size: 1, signed: false, repr: nil}
	fTShort     = &fieldType{id: 3, name: "SHORT", size: 2, signed: false, repr: nil}
	fTLong      = &fieldType{id: 4, name: "LONG", size: 4, signed: false, repr: nil}
	fTRational  = &fieldType{id: 5, name: "RATIONAL", size: 8, signed: false, repr: nil}
	fTSByte     = &fieldType{id: 6, name: "SBYTE", size: 1, signed: true, repr: nil}
	fTUndefined = &fieldType{id: 7, name: "UNDEFINED", size: 1, signed: false, repr: nil}
	fTSShort    = &fieldType{id: 8, name: "SSHORT", size: 2, signed: true, repr: nil}
	fTSLong     = &fieldType{id: 9, name: "SLONG", size: 4, signed: true, repr: nil}
	fTSRational = &fieldType{id: 10, name: "SRATIONAL", size: 8, signed: true, repr: nil}
	fTFloat     = &fieldType{id: 11, name: "FLOAT", size: 4, signed: true, repr: nil}
	fTDouble    = &fieldType{id: 12, name: "DOUBLE", size: 8, signed: true, repr: nil}
	fTIFD       = &fieldType{id: 13, name: "IFD", size: 4, signed: false, repr: nil}     // TODO: Double check parameters.
	fTUnicode   = &fieldType{id: 14, name: "UNICODE", size: 4, signed: false, repr: nil} // TODO: Double check parameters.
	fTComplex   = &fieldType{id: 15, name: "COMPLEX", size: 8, signed: true, repr: nil}  // TODO: Double check parameters.
	fTLong8     = &fieldType{id: 16, name: "LONG8", size: 8, signed: false, repr: nil}
	fTSLong8    = &fieldType{id: 17, name: "SLONG8", size: 8, signed: true, repr: nil}
	fTIFD8      = &fieldType{id: 18, name: "IFD8", size: 8, signed: false, repr: nil}
)
