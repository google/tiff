package tiff

import (
	"encoding/binary"
)

type FieldValue interface {
	Order() binary.ByteOrder
	Bytes() []byte
}

type fieldValue struct {
	order binary.ByteOrder
	value []byte
}

func (fv *fieldValue) Order() binary.ByteOrder {
	return fv.order
}

func (fv *fieldValue) Bytes() []byte {
	return fv.value
}

type Field interface {
	Tag() Tag
	Type() FieldType
	Count() uint32
	Value() FieldValue
}

type field struct {
	entry Entry

	// If the offset from entry is actually a value, then the bytes will be
	// stored here and offset should be considered 0.  Otherwise, the offset
	// will indicate the location in the file where the values can be found.
	// Then value will be used to hold the bytes associated with those values.
	value FieldValue

	// etSet is the FieldTypeSet that can be used to look up the FieldType
	// that corresponds to the typeId of this entry.  If etSet is nil, the
	// default set DefaultFieldTypes is used instead.
	etSet FieldTypeSet

	// tSet is the TagSet that can be used to look up the Tag that
	// corresponds to the result of entry.TagId().
	tSet TagSet
}

func (f *field) Tag() Tag {
	if f.tSet == nil {
		return DefaultTags.GetTag(f.entry.TagId())
	}
	return f.tSet.GetTag(f.entry.TagId())
}

func (f *field) Type() FieldType {
	if f.etSet == nil {
		return DefaultFieldTypes.GetType(f.entry.TypeId())
	}
	return f.etSet.GetType(f.entry.TypeId())
}

func (f *field) Count() uint32 {
	return f.entry.Count()
}

func (f *field) Value() FieldValue {
	return f.value
}

func parseField(br *bReader) (out Field, err error) {
	f := new(field)
	if f.entry, err = parseEntry(br); err != nil {
		return
	}
	// TODO: Implement grabbing the value.  For now, just use the bytes from
	// the ValueOffset.
	valOff := f.entry.ValueOffset()
	f.value = &fieldValue{
		order: br.order,
		value: valOff[:],
	}
	return f, nil
}

type Field8 interface {
	Tag() Tag
	Type() FieldType
	Count() uint64
	Value() FieldValue
}

type field8 struct {
	entry Entry8

	// If the offset entry is actually a value, then the bytes will be
	// stored here and offset should be set to 0.  Otherwise, the offset
	// will indicate the location in the file where the values can be found.
	// Then value will hold the bytes associated with those values.
	value FieldValue

	// etSet is the FieldTypeSet that can be used to look up the FieldType
	// that corresponds to the typeId of this entry.  If etSet is nil, the
	// default set DefaultFieldTypes is used instead.
	etSet FieldTypeSet

	// tSet is the TagSet that can be used to look up the Tag that
	// corresponds to the result of entry.TagId().
	tSet TagSet
}

func (f8 *field8) Tag() Tag {
	if f8.tSet == nil {
		return DefaultTags.GetTag(f8.entry.TagId())
	}
	return f8.tSet.GetTag(f8.entry.TagId())
}

func (f8 *field8) Type() FieldType {
	if f8.etSet == nil {
		return DefaultFieldTypes.GetType(f8.entry.TypeId())
	}
	return f8.etSet.GetType(f8.entry.TypeId())
}

func (f8 *field8) Count() uint64 {
	return f8.entry.Count()
}

func (f8 *field8) Value() FieldValue {
	return f8.value
}

func parseField8(br *bReader) (out Field8, err error) {
	f := new(field8)
	if f.entry, err = parseEntry8(br); err != nil {
		return
	}
	// TODO: Implement grabbing the value.  For now, just use the bytes from
	// the ValueOffset.
	valOff := f.entry.ValueOffset()
	f.value = &fieldValue{
		order: br.order,
		value: valOff[:],
	}
	return f, nil
}
