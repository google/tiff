package tiff

import (
	"encoding/binary"
	"fmt"
	"sync"
)

// FieldTypeSet represents a set of field types that may be in use within a file
// that uses a TIFF file structure.  This can be customized for custom file
// formats and private IFDs.
type FieldTypeSet interface {
	Register(ft FieldType) error
	GetType(id uint16) FieldType
}

type fieldTypeSet struct {
	mu    sync.Mutex
	types map[uint16]FieldType
}

func (fts *fieldTypeSet) Register(ft FieldType) error {
	fts.mu.Lock()
	defer fts.mu.Unlock()
	id := ft.Id()
	current, ok := fts.types[id]
	if ok {
		// If there is a need to overwrite a field type to use a different name or
		// size, then that probably belongs in a custom field type set that the
		// user can implement themselves for use in private IFDs.  We do not
		// want to override any of the size or name settings for any of the
		// default field types defined in this package.
		if current.Name() != ft.Name() {
			return fmt.Errorf("tiff: field type registration failure for id %d, name mismatch (current: %q, new: %q)", id, current.Name(), ft.Name())
		}
		if current.Size() != ft.Size() {
			return fmt.Errorf("tiff: field type registration failure for id %d, size mismatch (current: %d, new: %d)", id, current.Size(), ft.Size())
		}
	}

	// At this point, we are probably registering a new field type or a field
	// type that already exists with the same parameters.  We allow users to
	// register the same type over again in case they want to set and use a
	// different representation func.
	fts.types[id] = ft
	return nil
}

func (fts *fieldTypeSet) GetType(id uint16) FieldType {
	fts.mu.Lock()
	defer fts.mu.Unlock()
	if et, ok := fts.types[id]; ok {
		return et
	}
	return &fieldType{
		id:   id,
		name: fmt.Sprintf("UnregisteredFieldType_%d", id),
		repr: func([]byte, binary.ByteOrder) string { return "UNKNOWN" },
	}
}

// Note: We could create key and value pairs in the map by doing:
//     fTByte.Id(): fTByte,
// However, we know these values to be accurate with the implementations defined
// above.  Using the constant values reads better.  Any additions to this set
// should double check the values used above in the struct definition and here
// in the map key.
var defFieldTypes = &fieldTypeSet{
	types: map[uint16]FieldType{
		1:  fTByte,
		2:  fTASCII,
		3:  fTShort,
		4:  fTLong,
		5:  fTRational,
		6:  fTSByte,
		7:  fTUndefined,
		8:  fTSShort,
		9:  fTSLong,
		10: fTSRational,
		11: fTFloat,
		12: fTDouble,
		13: fTIFD,
		14: fTUnicode,
		15: fTComplex,
		16: fTLong8,
		17: fTSLong8,
		18: fTIFD8,
	},
}

// DefaultFieldTypes is the default set of field types supported by this
// package.  A user is free to create their own FieldTypeSet from which to
// support extended functionality or to provide a substitute representation for
// known types.  Most users will be fine with the default set defined here.
var DefaultFieldTypes FieldTypeSet = defFieldTypes

// RegisterFieldType allows a user to extend the default set of field types that may
// be in use in custom file formats that use the TIFF file structure.  Field types
// added via RegisterFieldType are added to the built-in default set assuming they
// do not conflict with existing field parameters.
func RegisterFieldType(ft FieldType) error {
	return DefaultFieldTypes.Register(ft)
}
