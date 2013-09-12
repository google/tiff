package tiff

/* Entry/Tag structure
For IFD/Tag:
	Each 12-byte IFD entry has the following format:
	Bytes 0-1:  The Tag that identifies the field.
	Bytes 2-3:  The field Type.
	Bytes 4-7:  The number of values, Count of the indicated Type.
	Bytes 8-11: The Value Offset, the file offset (in bytes) of the Value
	            for the field. The Value is expected to begin on a word
	            boundary; the corresponding Value Offset will thus be an
	            even number. This file offset may point anywhere in the
	            file, even after the image data.

For IFD8/Tag8:
	Each 20-byte IFD entry has the following format:
	Bytes 0-1:   The Tag that identifies the field.
	Bytes 2-3:   The field Type.
	Bytes 4-11:  The number of values, Count of the indicated Type.
	Bytes 12-19: The Value Offset, the file offset (in bytes) of the Value
	             for the field. The Value is expected to begin on a word
	             boundary; the corresponding Value Offset will thus be an
	             even number. This file offset may point anywhere in the
	             file, even after the image data.

*/

type Tag interface {
	Id() uint16
	Name() string
	TypeId() uint16
	Type() TagType
	Count() uint32
	Offset() uint32
	Value() []byte
}

type Tag8 interface {
	Id() uint16
	Name() string
	TypeId() uint16
	Type() TagType
	Count() uint64
	Offset() uint64
	Value() []byte
}
