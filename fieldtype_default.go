package tiff

import (
	"encoding/binary"
	"fmt"
	"math"
	"math/big"
)

/* Field type definitions
From [TIFF6]:
	1-12: Field types 1 - 12 are described in [TIFF6].

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
	13-15: Field Type IDs 13 - 15 are mentioned in [BIGTIFFDESIGN], but only
		to explain why values 13 - 15 were skipped when identifying new
		Field Types for BigTIFF. These are meant to be used with regular
		TIFF, but were apparently not properly documented prior to the
		BigTIFF design discussion.

	13 = IFD	?? 32-bit unsigned integer offset value ??
	14 = UNICODE	??
	15 = COMPLEX	??
*/

// Default set of Field types.  These are exported for others to use in
// registering custom tags.
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

// These functions provide string representations of values based on field types.
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

// DefaultFieldTypeSet is the default set of field types supported by this
// package.  A user is free to create their own FieldTypeSet from which to
// support extended functionality or to provide a substitute representation for
// known types.  Most users will be fine with the default set defined here.
var DefaultFieldTypeSet = NewFieldTypeSet("Default")

func init() {
	DefaultFieldTypeSet.Register(FTByte)
	DefaultFieldTypeSet.Register(FTAscii)
	DefaultFieldTypeSet.Register(FTShort)
	DefaultFieldTypeSet.Register(FTLong)
	DefaultFieldTypeSet.Register(FTRational)
	DefaultFieldTypeSet.Register(FTSByte)
	DefaultFieldTypeSet.Register(FTUndefined)
	DefaultFieldTypeSet.Register(FTSShort)
	DefaultFieldTypeSet.Register(FTSLong)
	DefaultFieldTypeSet.Register(FTSRational)
	DefaultFieldTypeSet.Register(FTFloat)
	DefaultFieldTypeSet.Register(FTDouble)
	DefaultFieldTypeSet.Register(FTIFD)
	DefaultFieldTypeSet.Register(FTUnicode)
	DefaultFieldTypeSet.Register(FTComplex)

	// Prevent further registration in the DefaultFieldTypeSet.  Others should
	// add to the DefaultFieldTypeSpace instead of the core set.
	DefaultFieldTypeSet.Lock()

	DefaultFieldTypeSpace.RegisterFieldTypeSet(DefaultFieldTypeSet)
}
