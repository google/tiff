package tiff

import (
	"fmt"
)

type Tag interface {
	ID() uint16
	Name() string
	ValidFieldTypes() []FieldType
}

func NewTag(id uint16, name string, validFTs []FieldType) Tag {
	return &tag{id: id, name: name, validFTs: validFTs}
}

type tag struct {
	id       uint16
	name     string
	validFTs []FieldType
}

func (t *tag) ID() uint16 {
	return t.id
}

func (t *tag) Name() string {
	if len(t.name) > 0 {
		return t.name
	}
	return fmt.Sprintf("UNNAMED_%d", t.id)
}

func (t *tag) ValidFieldTypes() []FieldType {
	return t.validFTs
}
