package tiff

import (
	"encoding/binary"
	"fmt"
	"sync"
)

/* Tag type definitions
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
	18 = IFD8	8-byte
*/

// A TagType represents all of the necessary pieces of information one needs to
// know about a tag type including a function that knows how to represent that
// type of data in an often human readable string.  Other string representation
// formats could be implemented (json, xml, etc).  Tag types themselves have no
// actual value inside a TIFF.  They are here to help an implementer or user
// understand their format.
type TagType interface {
	Id() uint16
	Name() string
	Size() uint64
	Signed() bool
	Repr() func([]byte, binary.ByteOrder) string
}

type tagType struct {
	id     uint16
	name   string
	size   uint64
	signed bool
	repr   func([]byte, binary.ByteOrder) string
}

func (tt *tagType) Id() uint16 {
	return tt.id
}

func (tt *tagType) Name() string {
	return tt.name
}

func (tt *tagType) Size() uint64 {
	return tt.size
}

func (tt *tagType) Signed() bool {
	return tt.signed
}

func (tt *tagType) Repr() func([]byte, binary.ByteOrder) string {
	return tt.repr
}

// TODO: Implement the basic representation func for each default tag type.

/* Default set of Tag types
1-12:  Tag types 1 - 12 are described in [TIFF6].
13-15: Tag Type IDs 13 - 15 are mentioned in [BIGTIFFDESIGN], but only to
       explain why values 13 - 15 were skipped when identifying new Tag Types
       for BigTIFF. These are meant to be used with regular TIFF, but were
       apparently not properly documented prior to the BigTIFF design
       discussion.
16-18: Tag Type IDs 16 - 18 were added for use with BigTIFF in [BIGTIFFDESIGN].
*/
var (
	// TTByte implements TagType with an Id of 1 and a composition of an
	// 8-bit unsigned integer.
	TTByte = &tagType{
		id:     1,
		name:   "BYTE",
		size:   1,
		signed: false,
		repr:   nil,
	}
	// TTASCII
	TTASCII = &tagType{
		id:     2,
		name:   "ASCII",
		size:   1,
		signed: false,
		repr:   nil,
	}
	// TTShort
	TTShort = &tagType{
		id:     3,
		name:   "SHORT",
		size:   2,
		signed: false,
		repr:   nil,
	}
	// TTLong
	TTLong = &tagType{
		id:     4,
		name:   "LONG",
		size:   4,
		signed: false,
		repr:   nil,
	}
	// TTRational
	TTRational = &tagType{
		id:     5,
		name:   "RATIONAL",
		size:   8,
		signed: false,
		repr:   nil,
	}
	// TTSByte
	TTSByte = &tagType{
		id:     6,
		name:   "SBYTE",
		size:   1,
		signed: true,
		repr:   nil,
	}
	// TTUndefined
	TTUndefined = &tagType{
		id:     7,
		name:   "UNDEFINED",
		size:   1,
		signed: false,
		repr:   nil,
	}
	// TTSShort
	TTSShort = &tagType{
		id:     8,
		name:   "SSHORT",
		size:   2,
		signed: true,
		repr:   nil,
	}
	// TTSLong
	TTSLong = &tagType{
		id:     9,
		name:   "SLONG",
		size:   4,
		signed: true,
		repr:   nil,
	}
	// TTSRational
	TTSRational = &tagType{
		id:     10,
		name:   "SRATIONAL",
		size:   8,
		signed: true,
		repr:   nil,
	}
	// TTFloat
	TTFloat = &tagType{
		id:     11,
		name:   "FLOAT",
		size:   4,
		signed: true,
		repr:   nil,
	}
	// TTDouble
	TTDouble = &tagType{
		id:     12,
		name:   "DOUBLE",
		size:   8,
		signed: true,
		repr:   nil,
	}
	// TTIFD
	TTIFD = &tagType{
		id:     13,
		name:   "IFD",
		size:   4,
		signed: false,
		repr:   nil,
	}
	// TTUnicode
	TTUnicode = &tagType{
		id:     14,
		name:   "UNICODE",
		size:   0, // TODO: Find out this size. int32?
		signed: false,
		repr:   nil,
	}
	// TTComplex
	TTComplex = &tagType{
		id:     15,
		name:   "COMPLEX",
		size:   8,    // 8 is a guess without proper documentation.
		signed: true, // true is a guess
		repr:   nil,
	}
	// TTLong8
	TTLong8 = &tagType{
		id:     16,
		name:   "LONG8",
		size:   8,
		signed: false,
		repr:   nil,
	}
	// TTSLong8
	TTSLong8 = &tagType{
		id:     17,
		name:   "SLONG8",
		size:   8,
		signed: true,
		repr:   nil,
	}
	// TTIFD8
	TTIFD8 = &tagType{
		id:     18,
		name:   "IFD8",
		size:   8,
		signed: false,
		repr:   nil,
	}
)

// TagTypeSet represents a set of tag types that may be in use within a file
// that uses a TIFF file structure.  This can be customized for custom file
// formats and private IFDs.
type TagTypeSet interface {
	Register(tt TagType) error
	GetType(id uint16) (TagType, error)
}

type tagTypeSet struct {
	mu    sync.Mutex
	types map[uint16]TagType
}

func (tts *tagTypeSet) Register(tt TagType) error {
	tts.mu.Lock()
	defer tts.mu.Unlock()
	id := tt.Id()
	current, ok := tts.types[id]
	if ok {
		// If there is a need to overwrite a tag type to use a different name or
		// size, then that probably belongs in a custom tag type set that the
		// user can implement themselves for use in private IFDs.  We do not
		// want to override any of the size or name settings for any of the
		// default tag types defined in this package.
		if current.Name() != tt.Name() {
			return fmt.Errorf("tiff: tag type registration failure for id %d, name mismatch (current: %q, new: %q)", id, current.Name(), tt.Name())
		}
		if current.Size() != tt.Size() {
			return fmt.Errorf("tiff: tag type registration failure for id %d, size mismatch (current: %d, new: %d)", id, current.Size(), tt.Size())
		}
	}

	// At this point, we are probably registering a new tag type or a tag
	// type that already exists with the same parameters.  We allow users to
	// register the same type over again in case they want to set and use a
	// different representation func.
	tts.types[id] = tt
	return nil
}

func (tts *tagTypeSet) GetType(id uint16) (TagType, error) {
	tts.mu.Lock()
	defer tts.mu.Unlock()
	tt, ok := tts.types[id]
	if !ok {
		return nil, fmt.Errorf("tiff: TagType id %d has not been registered", id)
	}
	return tt, nil
}

// Note: We could create key and value pairs in the map by doing:
//     TTByte.Id(): TTByte,
// However, we know these values to be accurate with the implementations defined
// above.  Using the constant values reads better.  Any additions to this set
// should double check the values used above in the struct definition and here
// in the map key.
var defTagTypes = &tagTypeSet{
	types: map[uint16]TagType{
		1:  TTByte,
		2:  TTASCII,
		3:  TTShort,
		4:  TTLong,
		5:  TTRational,
		6:  TTSByte,
		7:  TTUndefined,
		8:  TTSShort,
		9:  TTSLong,
		10: TTSRational,
		11: TTFloat,
		12: TTDouble,
		13: TTIFD,
		14: TTUnicode,
		15: TTComplex,
		16: TTLong8,
		17: TTSLong8,
		18: TTIFD8,
	},
}

// DefaultTagTypes is the default set of tag types supported by this package.  A
// user is free to create their own TagTypeSet from which to support extended
// functionality or to provide a substitute representation for known types.
// Most users will be fine with the default set defined here.
var DefaultTagTypes TagTypeSet = defTagTypes

// RegisterTagType allows a user to extend the default set of tag types that may
// be in use in custom file formats that use the TIFF file structure.  Tag types
// added via RegisterTagType are added to the built-in default set assuming they
// do not conflict with existing tag parameters.
func RegisterTagType(tt TagType) error {
	return defTagTypes.Register(tt)
}
