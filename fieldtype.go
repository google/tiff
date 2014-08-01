package tiff

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math"
	"math/big"
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
*/

// A FieldType represents all of the necessary pieces of information one needs to
// know about a field type including a function that knows how to represent that
// type of data in an often human readable string.  Other string representation
// formats could be implemented (json, xml, etc).  Field types themselves have no
// actual stored value inside a TIFF.  They are here to help an implementer or
// user understand their format.
type FieldType interface {
	ID() uint16
	Name() string
	Size() uint32
	Signed() bool
	Repr() func([]byte, binary.ByteOrder) string
}

func NewFieldType(id uint16, name string, size uint32, signed bool, repr func([]byte, binary.ByteOrder) string) FieldType {
	return &fieldType{id: id, name: name, size: size, signed: signed, repr: repr}
}

type fieldType struct {
	id     uint16
	name   string
	size   uint32
	signed bool
	repr   func([]byte, binary.ByteOrder) string
}

func (ft *fieldType) ID() uint16 {
	return ft.id
}

func (ft *fieldType) Name() string {
	return ft.name
}

func (ft *fieldType) Size() uint32 {
	return ft.size
}

func (ft *fieldType) Signed() bool {
	return ft.signed
}

func (ft *fieldType) Repr() func([]byte, binary.ByteOrder) string {
	return ft.repr
}

func (ft *fieldType) MarshalJSON() ([]byte, error) {
	tmp := struct {
		ID     uint16
		Name   string
		Size   uint32
		Signed bool
	}{
		ID:     ft.id,
		Name:   ft.name,
		Size:   ft.size,
		Signed: ft.signed,
	}
	return json.Marshal(tmp)
}

func reprByte(in []byte, bo binary.ByteOrder) string   { return fmt.Sprintf("%d", in[0]) }
func reprSByte(in []byte, bo binary.ByteOrder) string  { return fmt.Sprintf("%d", int8(in[0])) }
func reprASCII(in []byte, bo binary.ByteOrder) string  { return string(in) }
func reprShort(in []byte, bo binary.ByteOrder) string  { return fmt.Sprintf("%d", bo.Uint16(in)) }
func reprSShort(in []byte, bo binary.ByteOrder) string { return fmt.Sprintf("%d", int16(bo.Uint16(in))) }
func reprLong(in []byte, bo binary.ByteOrder) string   { return fmt.Sprintf("%d", bo.Uint32(in)) }
func reprSLong(in []byte, bo binary.ByteOrder) string  { return fmt.Sprintf("%d", int32(bo.Uint32(in))) }
func reprRational(in []byte, bo binary.ByteOrder) string {
	return big.NewRat(int64(bo.Uint32(in)), int64(bo.Uint32(in[4:]))).String()
}
func reprSRational(in []byte, bo binary.ByteOrder) string {
	return big.NewRat(int64(int32(bo.Uint32(in))), int64(int32(bo.Uint32(in[4:])))).String()
}
func reprFloat(in []byte, bo binary.ByteOrder) string {
	return fmt.Sprintf("%f", math.Float32frombits(bo.Uint32(in)))
}
func reprDouble(in []byte, bo binary.ByteOrder) string {
	return fmt.Sprintf("%f", math.Float64frombits(bo.Uint64(in)))
}

/*
Default set of Field types
  1-12:  Field types 1 - 12 are described in [TIFF6].
  13-15: Field Type IDs 13 - 15 are mentioned in [BIGTIFFDESIGN], but only to
         explain why values 13 - 15 were skipped when identifying new Field Types
         for BigTIFF. These are meant to be used with regular TIFF, but were
         apparently not properly documented prior to the BigTIFF design
         discussion.
*/
var (
	FTByte      = NewFieldType(1, "BYTE", 1, false, reprByte)
	FTAscii     = NewFieldType(2, "ASCII", 1, false, reprASCII)
	FTShort     = NewFieldType(3, "SHORT", 2, false, reprShort)
	FTLong      = NewFieldType(4, "LONG", 4, false, reprLong)
	FTRational  = NewFieldType(5, "RATIONAL", 8, false, reprRational)
	FTSByte     = NewFieldType(6, "SBYTE", 1, true, reprSByte)
	FTUndefined = NewFieldType(7, "UNDEFINED", 1, false, reprByte)
	FTSShort    = NewFieldType(8, "SSHORT", 2, true, reprSShort)
	FTSLong     = NewFieldType(9, "SLONG", 4, true, reprSLong)
	FTSRational = NewFieldType(10, "SRATIONAL", 8, true, reprSRational)
	FTFloat     = NewFieldType(11, "FLOAT", 4, true, reprFloat)
	FTDouble    = NewFieldType(12, "DOUBLE", 8, true, reprDouble)

	// TODO: The following 3 field types are not well defined.  Double check
	// their parameters.

	FTIFD     = NewFieldType(13, "IFD", 4, false, nil)
	FTUnicode = NewFieldType(14, "UNICODE", 4, false, nil)
	FTComplex = NewFieldType(15, "COMPLEX", 8, true, nil)
)
