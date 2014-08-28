package tiff

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
)

// A FieldType represents all of the necessary pieces of information one needs to
// know about a field type including a function that knows how to represent that
// type of data in an often human readable string.  Other string representation
// formats could be implemented (json, xml, etc).  Field types themselves have no
// actual stored value inside a TIFF.  They are here to help an implementer or
// user understand their format.
type FieldType interface {
	ID() uint16
	Name() string
	Size() uint64
	Signed() bool
	Repr() func([]byte, binary.ByteOrder) string
}

func NewFieldType(id uint16, name string, size uint64, signed bool, repr func([]byte, binary.ByteOrder) string) FieldType {
	return &fieldType{id: id, name: name, size: size, signed: signed, repr: repr}
}

type fieldType struct {
	id     uint16
	name   string
	size   uint64
	signed bool
	repr   func([]byte, binary.ByteOrder) string
}

func (ft *fieldType) ID() uint16 {
	return ft.id
}

func (ft *fieldType) Name() string {
	if len(ft.name) > 0 {
		return ft.name
	}
	return fmt.Sprintf("UNNAMED_FIELD_TYPE_%d", ft.id)
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

func (ft *fieldType) MarshalJSON() ([]byte, error) {
	tmp := struct {
		ID     uint16
		Name   string
		Size   uint64
		Signed bool
	}{
		ID:     ft.id,
		Name:   ft.name,
		Size:   ft.size,
		Signed: ft.signed,
	}
	return json.Marshal(tmp)
}
