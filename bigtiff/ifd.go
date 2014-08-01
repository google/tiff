package bigtiff

import "github.com/jonathanpittman/tiff"

// IFD represents the data structure of an IFD in a BigTIFF.
type IFD interface {
	NumEntries() uint64
	Fields() []Field
	NextOffset() uint64
}

type imageFileDirectory struct {
	numEntries uint64
	fields     []Field
	nextOffset uint64
}

func (ifd *imageFileDirectory) NumEntries() uint64 {
	return ifd.numEntries
}

func (ifd *imageFileDirectory) Fields() []Field {
	return ifd.fields
}

func (ifd *imageFileDirectory) NextOffset() uint64 {
	return ifd.nextOffset
}

func ParseIFD(br tiff.BReader, offset uint64, tsg tiff.TagSpace, ftsp tiff.FieldTypeSpace) (out IFD, err error) {
	ifd := new(imageFileDirectory)
	br.Seek(int64(offset), 0) // TODO: This may be wrong.  Need uint64 capacity?
	if err = br.BRead(&ifd.numEntries); err != nil {
		return
	}
	for i := uint64(0); i < ifd.numEntries; i++ {
		var f Field
		if f, err = ParseField(br, tsg, ftsp); err != nil {
			return
		}
		ifd.fields = append(ifd.fields, f)
	}
	if err = br.BRead(&ifd.nextOffset); err != nil {
		return
	}
	return ifd, nil
}
