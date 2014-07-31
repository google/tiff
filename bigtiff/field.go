package bigtiff

import (
	"encoding/binary"
	"encoding/json"

	"github.com/jonathanpittman/tiff"
)

// A Field is primarily comprised of an Entry.  If the Entry's value is actually
// an offset to the data that the entry describes, then the Field will contain
// both the offset and the data that the offset points to in the file.

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

func (fv *fieldValue) MarshalJSON() ([]byte, error) {
	tmp := struct {
		Bytes []byte
	}{
		Bytes: fv.value,
	}
	return json.Marshal(tmp)
}

// Field8 represents a field in an IFD8 for a BigTIFF file.
type Field8 interface {
	Tag() tiff.Tag
	Type() tiff.FieldType
	Count() uint64
	Offset() uint64
	Value() tiff.FieldValue
}

type field8 struct {
	entry Entry8

	// If the offset entry is actually a value, then the bytes will be
	// stored here and offset should be set to 0.  Otherwise, the offset
	// will indicate the location in the file where the values can be found.
	// Then value will hold the bytes associated with those values.
	value tiff.FieldValue

	// fts is the FieldTypeSet that can be used to look up the FieldType
	// that corresponds to the typeId of this entry.  If fts is nil, the
	// default set DefaultFieldTypes is used instead.
	fts tiff.FieldTypeSet

	// tsg is the TagSpace that can be used to look up the Tag that
	// corresponds to the result of entry.TagID().
	tsg tiff.TagSpace
}

func (f8 *field8) Tag() tiff.Tag {
	if f8.tsg == nil {
		return tiff.DefaultTagSpace.GetTag(f8.entry.TagID())
	}
	return f8.tsg.GetTag(f8.entry.TagID())
}

func (f8 *field8) Type() tiff.FieldType {
	if f8.fts == nil {
		return tiff.DefaultFieldTypes.GetType(f8.entry.TypeID())
	}
	return f8.fts.GetType(f8.entry.TypeID())
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

func (f8 *field8) Value() tiff.FieldValue {
	return f8.value
}

func (f8 *field8) MarshalJSON() ([]byte, error) {
	return nil, nil
}

func ParseField8(br tiff.BReader, tsg tiff.TagSpace, fts tiff.FieldTypeSet) (out Field8, err error) {
	if fts == nil {
		fts = tiff.DefaultFieldTypes
	}
	if tsg == nil {
		tsg = tiff.DefaultTagSpace
	}
	f := &field8{fts: fts, tsg: tsg}
	if f.entry, err = ParseEntry8(br); err != nil {
		return
	}
	// TODO: Implement grabbing the value.  For now, just use the bytes from
	// the ValueOffset.
	valOff := f.entry.ValueOffset()
	f.value = &fieldValue{
		order: br.Order(),
		value: valOff[:],
	}
	return f, nil
}
