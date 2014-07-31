package bigtiff

import "github.com/jonathanpittman/tiff"

// IFD8 represents the data structure of an IFD in a BigTIFF file.
type IFD8 interface {
	NumEntries() uint64
	Fields() []Field8
	NextOffset() uint64
}

type imageFileDirectory8 struct {
	numEntries uint64
	fields     []Field8
	nextOffset uint64
}

func (ifd8 *imageFileDirectory8) NumEntries() uint64 {
	return ifd8.numEntries
}

func (ifd8 *imageFileDirectory8) Fields() []Field8 {
	return ifd8.fields
}

func (ifd8 *imageFileDirectory8) NextOffset() uint64 {
	return ifd8.nextOffset
}

func ParseIFD8(br tiff.BReader, offset uint64, tsg tiff.TagSpace, fts tiff.FieldTypeSet) (out IFD8, err error) {
	ifd := new(imageFileDirectory8)
	br.Seek(int64(offset), 0) // TODO: This is wrong.  Need uint64.  Use big.Int?
	if err = br.BRead(&ifd.numEntries); err != nil {
		return
	}
	for i := uint64(0); i < ifd.numEntries; i++ {
		var f Field8
		if f, err = ParseField8(br, tsg, fts); err != nil {
			return
		}
		ifd.fields = append(ifd.fields, f)
	}
	if err = br.BRead(&ifd.nextOffset); err != nil {
		return
	}
	return ifd, nil
}
