package tiff

import (
	"encoding/binary"
	"fmt"
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
	Offset() uint32
	Value() FieldValue
}

type field struct {
	entry Entry

	// If the offset from entry is actually a value, then the bytes will be
	// stored here and offset should be considered 0.  Otherwise, the offset
	// will indicate the location in the file where the values can be found.
	// Then value will be used to hold the bytes associated with those values.
	value FieldValue

	// fts is the FieldTypeSet that can be used to look up the FieldType
	// that corresponds to the typeId of this entry.  If fts is nil, the
	// default set DefaultFieldTypes is used instead.
	fts FieldTypeSet

	// tSet is the TagSetGroup that can be used to look up the Tag that
	// corresponds to the result of entry.TagId().
	tsg TagSetGroup
}

func (f *field) Tag() Tag {
	if f.tsg == nil {
		return DefaultTagSetGroup.GetTag(f.entry.TagId())
	}
	return f.tsg.GetTag(f.entry.TagId())
}

func (f *field) Type() FieldType {
	if f.fts == nil {
		return DefaultFieldTypes.GetType(f.entry.TypeId())
	}
	return f.fts.GetType(f.entry.TypeId())
}

func (f *field) Count() uint32 {
	return f.entry.Count()
}

func (f *field) Offset() uint32 {
	if f.Type().Size()*f.Count() <= 4 {
		return 0
	}
	offsetBytes := f.entry.ValueOffset()
	return f.Value().Order().Uint32(offsetBytes[:])
}

func (f *field) Value() FieldValue {
	return f.value
}

func (f *field) String() string {
	var (
		theTSet  TagSetGroup  = f.tsg
		theFTSet FieldTypeSet = f.fts
	)
	if f.tsg == nil {
		theTSet = DefaultTagSetGroup
	}
	if f.fts == nil {
		theFTSet = DefaultFieldTypes
	}
	return fmt.Sprintf("<Tag: %v, Type: %v, Count: %d, Offset: %d, Value: %v, FieldTypeSet: %q, TagSetGroup: %q>",
		f.Tag().Name(), f.Type().Name(), f.Count(), f.Offset(), f.Value().Bytes(), theFTSet.Name(), theTSet.Name())
}

func parseField(br *bReader) (out Field, err error) {
	f := new(field)
	if f.entry, err = parseEntry(br); err != nil {
		return
	}
	fv := &fieldValue{order: br.order}
	valSize := int64(f.Count()) * int64(f.Type().Size())
	valOffBytes := f.entry.ValueOffset()
	if valSize > 4 {
		fv.value = make([]byte, valSize)
		offset := int64(br.order.Uint32(valOffBytes[:]))
		if err = br.ReadSection(offset, valSize, &fv.value); err != nil {
			return
		}
	} else {
		fv.value = valOffBytes[:]
	}
	f.value = fv
	return f, nil
}

type Field8 interface {
	Tag() Tag
	Type() FieldType
	Count() uint64
	Offset() uint64
	Value() FieldValue
}

type field8 struct {
	entry Entry8

	// If the offset entry is actually a value, then the bytes will be
	// stored here and offset should be set to 0.  Otherwise, the offset
	// will indicate the location in the file where the values can be found.
	// Then value will hold the bytes associated with those values.
	value FieldValue

	// fts is the FieldTypeSet that can be used to look up the FieldType
	// that corresponds to the typeId of this entry.  If fts is nil, the
	// default set DefaultFieldTypes is used instead.
	fts FieldTypeSet

	// tsg is the TagSetGroup that can be used to look up the Tag that
	// corresponds to the result of entry.TagId().
	tsg TagSetGroup
}

func (f8 *field8) Tag() Tag {
	if f8.tsg == nil {
		return DefaultTagSetGroup.GetTag(f8.entry.TagId())
	}
	return f8.tsg.GetTag(f8.entry.TagId())
}

func (f8 *field8) Type() FieldType {
	if f8.fts == nil {
		return DefaultFieldTypes.GetType(f8.entry.TypeId())
	}
	return f8.fts.GetType(f8.entry.TypeId())
}

func (f8 *field8) Count() uint64 {
	return f8.entry.Count()
}

func (f8 *field8) Offset() uint64 {
	if uint64(f8.Type().Size())*f8.Count() <= 8 {
		return 0
	}
	offsetBytes := f8.entry.ValueOffset()
	return f8.Value().Order().Uint64(offsetBytes[:])
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
