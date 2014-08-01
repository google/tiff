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

// Field represents a field in an IFD for a BigTIFF file.
type Field interface {
	Tag() tiff.Tag
	Type() tiff.FieldType
	Count() uint64
	Offset() uint64
	Value() tiff.FieldValue
}

type field struct {
	entry Entry

	// If the offset entry is actually a value, then the bytes will be
	// stored here and offset should be set to 0.  Otherwise, the offset
	// will indicate the location in the file where the values can be found.
	// Then value will hold the bytes associated with those values.
	value tiff.FieldValue

	// ftsp is the FieldTypeSpace that can be used to look up the FieldType
	// that corresponds to the typeId of this entry.  If ftsp is nil, the
	// default set tiff.DefaultFieldTypeSpace is used instead.
	ftsp tiff.FieldTypeSpace

	// tsp is the TagSpace that can be used to look up the Tag that
	// corresponds to the result of entry.TagID().
	tsp tiff.TagSpace
}

func (f *field) Tag() tiff.Tag {
	if f.tsp == nil {
		return tiff.DefaultTagSpace.GetTag(f.entry.TagID())
	}
	return f.tsp.GetTag(f.entry.TagID())
}

func (f *field) Type() tiff.FieldType {
	if f.ftsp == nil {
		return tiff.DefaultFieldTypeSpace.GetFieldType(f.entry.TypeID())
	}
	return f.ftsp.GetFieldType(f.entry.TypeID())
}

func (f *field) Count() uint64 {
	return f.entry.Count()
}

func (f *field) Offset() uint64 {
	if uint64(f.Type().Size())*f.Count() <= 8 {
		return 0
	}
	offsetBytes := f.entry.ValueOffset()
	return f.Value().Order().Uint64(offsetBytes[:])
}

func (f *field) Value() tiff.FieldValue {
	return f.value
}

func (f *field) MarshalJSON() ([]byte, error) {
	return nil, nil
}

func ParseField(br tiff.BReader, tsp tiff.TagSpace, ftsp tiff.FieldTypeSpace) (out Field, err error) {
	if ftsp == nil {
		ftsp = tiff.DefaultFieldTypeSpace
	}
	if tsp == nil {
		tsp = tiff.DefaultTagSpace
	}
	f := &field{ftsp: ftsp, tsp: tsp}
	if f.entry, err = ParseEntry(br); err != nil {
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
