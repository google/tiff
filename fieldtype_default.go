package tiff

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
	"math/big"
	"reflect"
)

/* FieldTypeRepr */
// These functions provide string representations of values based on field types.
func reprByte(in []byte, bo binary.ByteOrder) string   { return fmt.Sprintf("%d", in[0]) }
func reprSByte(in []byte, bo binary.ByteOrder) string  { return fmt.Sprintf("%d", int8(in[0])) }
func reprASCII(in []byte, bo binary.ByteOrder) string  { return string(in) }
func reprShort(in []byte, bo binary.ByteOrder) string  { return fmt.Sprintf("%d", bo.Uint16(in)) }
func reprSShort(in []byte, bo binary.ByteOrder) string { return fmt.Sprintf("%d", int16(bo.Uint16(in))) }
func reprLong(in []byte, bo binary.ByteOrder) string   { return fmt.Sprintf("%d", bo.Uint32(in)) }
func reprSLong(in []byte, bo binary.ByteOrder) string  { return fmt.Sprintf("%d", int32(bo.Uint32(in))) }
func reprRational(in []byte, bo binary.ByteOrder) string {
	// Print the representation directly to prevent panics from divide by
	// zero errors when using big.NewRat().
	return fmt.Sprintf("%d/%d", int64(bo.Uint32(in)), int64(bo.Uint32(in[4:])))
}
func reprSRational(in []byte, bo binary.ByteOrder) string {
	// Print the representation directly to prevent panics from divide by
	// zero errors when using big.NewRat()
	return fmt.Sprintf("%d/%d", int64(int32(bo.Uint32(in))), int64(int32(bo.Uint32(in[4:]))))
}
func reprFloat(in []byte, bo binary.ByteOrder) string {
	return fmt.Sprintf("%f", math.Float32frombits(bo.Uint32(in)))
}
func reprDouble(in []byte, bo binary.ByteOrder) string {
	return fmt.Sprintf("%f", math.Float64frombits(bo.Uint64(in)))
}

/* FieldTypeValuer */
func rvalByte(in []byte, bo binary.ByteOrder) reflect.Value  { return reflect.ValueOf(in[0]) }
func rvalSByte(in []byte, bo binary.ByteOrder) reflect.Value { return reflect.ValueOf(int8(in[0])) }
func rvalASCII(in []byte, bo binary.ByteOrder) reflect.Value {
	return reflect.ValueOf(string(bytes.TrimRight(in, "\x00")))
}
func rvalShort(in []byte, bo binary.ByteOrder) reflect.Value { return reflect.ValueOf(bo.Uint16(in)) }
func rvalSShort(in []byte, bo binary.ByteOrder) reflect.Value {
	return reflect.ValueOf(int16(bo.Uint16(in)))
}
func rvalLong(in []byte, bo binary.ByteOrder) reflect.Value { return reflect.ValueOf(bo.Uint32(in)) }
func rvalSLong(in []byte, bo binary.ByteOrder) reflect.Value {
	return reflect.ValueOf(int32(bo.Uint32(in)))
}
func rvalRational(in []byte, bo binary.ByteOrder) reflect.Value {
	denom := int64(bo.Uint32(in[4:]))
	if denom == 0 {
		// Prevent panics due to poorly written Rational fields with a
		// denominator of 0.
		return reflect.New(reflect.TypeOf(big.Rat{}))
	}
	numer := int64(bo.Uint32(in))
	return reflect.ValueOf(big.NewRat(numer, denom))
}
func rvalSRational(in []byte, bo binary.ByteOrder) reflect.Value {
	denom := int64(int32(bo.Uint32(in[4:])))
	if denom == 0 {
		// Prevent panics due to poorly written Rational fields with a
		// denominator of 0.  Their usable value would likely be 0.
		return reflect.New(reflect.TypeOf(big.Rat{}))
	}
	numer := int64(int32(bo.Uint32(in)))
	return reflect.ValueOf(big.NewRat(numer, denom))
}
func rvalFloat(in []byte, bo binary.ByteOrder) reflect.Value {
	return reflect.ValueOf(math.Float32frombits(bo.Uint32(in)))
}
func rvalDouble(in []byte, bo binary.ByteOrder) reflect.Value {
	return reflect.ValueOf(math.Float64frombits(bo.Uint64(in)))
}

/* reflect.Type */
var (
	typByte   = reflect.TypeOf(byte(0))         // BYTE, UNDEFINED
	typString = reflect.TypeOf(string(""))      // ASCII
	typU16    = reflect.TypeOf(uint16(0))       // SHORT
	typU32    = reflect.TypeOf(uint32(0))       // LONG, IFD
	typBigRat = reflect.TypeOf((*big.Rat)(nil)) // RATIONAL, SRATIONAL
	typI8     = reflect.TypeOf(int8(0))         // SBYTE
	typI16    = reflect.TypeOf(int16(0))        // SSHORT
	typI32    = reflect.TypeOf(int32(0))        // SLONG
	typF32    = reflect.TypeOf(float32(0))      // FLOAT
	typF64    = reflect.TypeOf(float64(0))      // DOUBLE
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
	FTByte      = NewFieldType(1, "BYTE", 1, false, reprByte, rvalByte, typByte)
	FTAscii     = NewFieldType(2, "ASCII", 1, false, reprASCII, rvalASCII, typString)
	FTShort     = NewFieldType(3, "SHORT", 2, false, reprShort, rvalShort, typU16)
	FTLong      = NewFieldType(4, "LONG", 4, false, reprLong, rvalLong, typU32)
	FTRational  = NewFieldType(5, "RATIONAL", 8, false, reprRational, rvalRational, typBigRat)
	FTSByte     = NewFieldType(6, "SBYTE", 1, true, reprSByte, rvalSByte, typI8)
	FTUndefined = NewFieldType(7, "UNDEFINED", 1, false, reprByte, rvalByte, typByte)
	FTSShort    = NewFieldType(8, "SSHORT", 2, true, reprSShort, rvalSShort, typI16)
	FTSLong     = NewFieldType(9, "SLONG", 4, true, reprSLong, rvalSLong, typI32)
	FTSRational = NewFieldType(10, "SRATIONAL", 8, true, reprSRational, rvalSRational, typBigRat)
	FTFloat     = NewFieldType(11, "FLOAT", 4, true, reprFloat, rvalFloat, typF32)
	FTDouble    = NewFieldType(12, "DOUBLE", 8, true, reprDouble, rvalDouble, typF64)
	FTIFD       = NewFieldType(13, "IFD", 4, false, reprLong, rvalLong, typU32)

	// TODO: These two are not complete.  Get the details and finish them.
	// For now, represent them as bytes until a proper representation is known.
	FTUnicode = NewFieldType(14, "UNICODE", 1, false, reprByte, rvalByte, typByte)
	FTComplex = NewFieldType(15, "COMPLEX", 1, true, reprByte, rvalByte, typByte)
)

/*
Regarding UNICODE and COMPLEX field types:
  UNICODE: It is thought that it is intended for UNICODE to be represented by a
  32bit integer value.  It is NOT KNOWN whether this should be signed or
  unsigned.  In go terms, this would either be an int32 or uint32.  The size
  would then be 4 instead of 1.  In tiff terms this would mirror either a LONG
  or an SLONG.

  COMPLEX: It is thought that it is intended for COMPLEX to be represented by
  two float32 values with the real part followed by the imaginary part in the
  byte stream.  In go terms this would mirror a complex64.  The size would then
  be 8 instead of 1.  In tiff terms this would be two FLOAT types in a similar
  way that a RATIONAL is two LONGs.
*/

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
