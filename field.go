package tiff

// A Field is primarily comprised of an Entry.  If the Entry's value is actually
// an offset to the data that the entry describes, then the Field will contain
// both the offset and the data that the offset points to in the file.

import (
	"encoding/binary"
	"encoding/json"
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

func (fv *fieldValue) MarshalJSON() ([]byte, error) {
	tmp := struct {
		Bytes []byte
	}{
		Bytes: fv.value,
	}
	return json.Marshal(tmp)
}

// Field represents a field in an IFD in a TIFF file.
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

	// tsg is the TagSpace that can be used to look up the Tag that
	// corresponds to the result of entry.TagId().
	tsg TagSpace
}

func (f *field) Tag() Tag {
	if f.tsg == nil {
		return DefaultTagSpace.GetTag(f.entry.TagId())
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
		theTSet  TagSpace     = f.tsg
		theFTSet FieldTypeSet = f.fts
	)
	if f.tsg == nil {
		theTSet = DefaultTagSpace
	}
	if f.fts == nil {
		theFTSet = DefaultFieldTypes
	}
	var valueRep string
	switch f.Type() {
	case FTAscii:
		valueRep = fmt.Sprintf("%q", f.Value().Bytes())
	default:
		// For general display purposes, don't show more than maxItems
		// amount of elements.  In this case, 10 is reasonable.  Beyond
		// 10 starts to line wrap in some cases like rationals.  If
		// we encounter that the value contains more than 10, we append
		// the ... to the end during string formatting to indicate that
		// there were more values, but they are not displayed here.
		const maxItems = 10
		buf := f.Value().Bytes()
		size := f.Type().Size()
		count := f.Count()
		if count > maxItems {
			count = maxItems
			buf = buf[:count*size]
		}
		vals := make([]string, 0, count)
		for len(buf) > 0 {
			if f.Type().Repr() != nil {
				vals = append(vals, f.Type().Repr()(buf[:size], f.Value().Order()))
			} else {
				vals = append(vals, fmt.Sprintf("%v", buf[:size]))
			}
			buf = buf[size:]
		}
		if count == 1 {
			valueRep = vals[0]
		} else if f.Count() > maxItems {
			valueRep = fmt.Sprintf("%v...", vals)
		} else {
			valueRep = fmt.Sprintf("%v", vals)
		}
	}
	return fmt.Sprintf("<Tag: %v, Type: %v, Count: %d, Offset: %d, Value: %s, FieldTypeSet: %q, TagSpace: %q>",
		f.Tag().Name(), f.Type().Name(), f.Count(), f.Offset(), valueRep, theFTSet.Name(), theTSet.Name())
}

func (f *field) MarshalJSON() ([]byte, error) {
	tmp := struct {
		E   Entry        `json:"Entry"`
		V   FieldValue   `json:"FieldValue"`
		FTS FieldTypeSet `json:"FieldTypeSet"`
		TSG TagSpace     `json:"TagSpace"`
	}{
		E:   f.entry,
		V:   f.value,
		FTS: f.fts,
		TSG: f.tsg,
	}
	return json.Marshal(tmp)
}

func ParseField(br BReader, tsg TagSpace, fts FieldTypeSet) (out Field, err error) {
	if fts == nil {
		fts = DefaultFieldTypes
	}
	if tsg == nil {
		tsg = DefaultTagSpace
	}
	f := &field{fts: fts, tsg: tsg}
	if f.entry, err = ParseEntry(br); err != nil {
		return
	}
	fv := &fieldValue{order: br.Order()}
	valSize := int64(f.Count()) * int64(f.Type().Size())
	valOffBytes := f.entry.ValueOffset()
	if valSize > 4 {
		fv.value = make([]byte, valSize)
		offset := int64(br.Order().Uint32(valOffBytes[:]))
		if err = br.BReadSection(&fv.value, offset, valSize); err != nil {
			return
		}
	} else {
		fv.value = valOffBytes[:]
	}
	f.value = fv
	return f, nil
}

// Field8 represents a field in an IFD8 for a BigTIFF file.
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

	// tsg is the TagSpace that can be used to look up the Tag that
	// corresponds to the result of entry.TagId().
	tsg TagSpace
}

func (f8 *field8) Tag() Tag {
	if f8.tsg == nil {
		return DefaultTagSpace.GetTag(f8.entry.TagId())
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

func (f8 *field8) MarshalJSON() ([]byte, error) {
	return nil, nil
}

func ParseField8(br BReader, tsg TagSpace, fts FieldTypeSet) (out Field8, err error) {
	if fts == nil {
		fts = DefaultFieldTypes
	}
	if tsg == nil {
		tsg = DefaultTagSpace
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
