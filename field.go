package tiff

type Field interface {
	Tag() Tag
	Type() FieldType
	Count() uint32
	Value() []byte
}

type field struct {
	entry Entry

	// If the offset entry is actually a value, then the bytes will be
	// stored here and offset should be set to 0.  Otherwise, the offset
	// will indicate the location in the file where the values can be found.
	// Then value will hold the bytes associated with those values.
	value []byte

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

func (f *field) Value() []byte {
	return f.value
}

type Field8 interface {
	Tag() Tag
	Type() FieldType
	Count() uint64
	Value() []byte
}

type field8 struct {
	entry Entry8

	// If the offset entry is actually a value, then the bytes will be
	// stored here and offset should be set to 0.  Otherwise, the offset
	// will indicate the location in the file where the values can be found.
	// Then value will hold the bytes associated with those values.
	value []byte

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

func (f8 *field8) Value() []byte {
	return f8.value
}
