package tiff

/*
Entry structure
  For IFD/Entry:
	Each 12-byte IFD entry has the following format:
	Bytes 0-1:  The Tag that identifies the entry.
	Bytes 2-3:  The entry Type.
	Bytes 4-7:  The number of values, Count of the indicated Type.
	Bytes 8-11: The Value Offset, the file offset (in bytes) of the Value
	            for the entry. The Value is expected to begin on a word
	            boundary; the corresponding Value Offset will thus be an
	            even number. This file offset may point anywhere in the
	            file, even after the image data.

  For IFD8/Entry8:
	Each 20-byte IFD entry has the following format:
	Bytes 0-1:   The Tag that identifies the entry.
	Bytes 2-3:   The entry Type.
	Bytes 4-11:  The number of values, Count of the indicated Type.
	Bytes 12-19: The Value Offset, the file offset (in bytes) of the Value
	             for the entry. The Value is expected to begin on a word
	             boundary; the corresponding Value Offset will thus be an
	             even number. This file offset may point anywhere in the
	             file, even after the image data.
*/

type Entry interface {
	TagId() uint16
	TypeId() uint16
	Count() uint32
	ValueOffset() [4]byte
}

// entry represents the data structure of an IFD entry.
type entry struct {
	tagId       uint16  // Bytes 0-1
	typeId      uint16  // Bytes 2-3
	count       uint32  // Bytes 4-7
	valueOffset [4]byte // Bytes 8-11
}

func (e *entry) TagId() uint16 {
	return e.tagId
}

func (e *entry) TypeId() uint16 {
	return e.typeId
}

func (e *entry) Count() uint32 {
	return e.count
}

func (e *entry) ValueOffset() [4]byte {
	return e.valueOffset
}

func parseEntry(br *bReader) (out Entry, err error) {
	e := new(entry)
	if err = br.Read(e); err != nil {
		return
	}
	return e, nil
}

type Entry8 interface {
	TagId() uint16
	TypeId() uint16
	Count() uint64
	ValueOffset() [8]byte
}

// entry8 represents the data structure of an IFD8 entry.
type entry8 struct {
	tagId       uint16  // Bytes 0-1
	typeId      uint16  // Bytes 2-3
	count       uint64  // Bytes 4-11
	valueOffset [8]byte // Bytes 12-19
}

func (e8 *entry8) TagId() uint16 {
	return e8.tagId
}

func (e8 *entry8) TypeId() uint16 {
	return e8.typeId
}

func (e8 *entry8) Count() uint64 {
	return e8.count
}

func (e8 *entry8) ValueOffset() [8]byte {
	return e8.valueOffset
}

func parseEntry8(br *bReader) (out Entry8, err error) {
	e := new(entry8)
	if err = br.Read(e); err != nil {
		return
	}
	return e, nil
}
