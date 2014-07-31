package tiff

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"sync"
)

// FieldTypeSet represents a set of field types that may be in use within a file
// that uses a TIFF file structure.  This can be customized for custom file
// formats and private IFDs.
type FieldTypeSet interface {
	Register(ft FieldType) error
	GetType(id uint16) FieldType
	Name() string
}

type fieldTypeSet struct {
	mu    sync.Mutex
	name  string
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

func (fts *fieldTypeSet) Name() string {
	return fts.name
}

func (fts *fieldTypeSet) MarshalJSON() ([]byte, error) {
	tmp := struct {
		Name string
	}{
		Name: fts.name,
	}
	return json.Marshal(tmp)
}

// Note: We could create key and value pairs in the map by doing:
//     fTByte.Id(): fTByte,
// However, we know these values to be accurate with the implementations defined
// above.  Using the constant values reads better.  Any additions to this set
// should double check the values used above in the struct definition and here
// in the map key.
var defFieldTypes = &fieldTypeSet{
	name: "Default",
	types: map[uint16]FieldType{
		1:  FTByte,
		2:  FTAscii,
		3:  FTShort,
		4:  FTLong,
		5:  FTRational,
		6:  FTSByte,
		7:  FTUndefined,
		8:  FTSShort,
		9:  FTSLong,
		10: FTSRational,
		11: FTFloat,
		12: FTDouble,
		13: FTIFD,
		14: FTUnicode,
		15: FTComplex,
		16: FTLong8,
		17: FTSLong8,
		18: FTIFD8,
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

var fieldTypeSets = struct {
	mu  sync.Mutex
	fts map[string]FieldTypeSet
}{
	fts: map[string]FieldTypeSet{DefaultFieldTypes.Name(): DefaultFieldTypes},
}

func RegisterFieldTypeSet(fts FieldTypeSet) error {
	fieldTypeSets.mu.Lock()
	defer fieldTypeSets.mu.Unlock()
	_, ok := fieldTypeSets.fts[fts.Name()]
	if ok {
		return fmt.Errorf("tiff: FieldTypeSet %q already registered.")
	}
	fieldTypeSets.fts[fts.Name()] = fts
	return nil
}

func GetFieldTypeSet(name string) (FieldTypeSet, error) {
	fieldTypeSets.mu.Lock()
	defer fieldTypeSets.mu.Unlock()
	return nil, nil
}
