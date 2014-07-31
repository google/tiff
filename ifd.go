package tiff

import "errors"

// IFD represents the data structure of an IFD in a TIFF File.
type IFD interface {
	NumEntries() uint16
	Fields() []Field
	NextOffset() uint32
}

type imageFileDirectory struct {
	numEntries uint16
	fields     []Field
	nextOffset uint32
}

func (ifd *imageFileDirectory) NumEntries() uint16 {
	return ifd.numEntries
}

func (ifd *imageFileDirectory) Fields() []Field {
	return ifd.fields
}

func (ifd *imageFileDirectory) NextOffset() uint32 {
	return ifd.nextOffset
}

func ParseIFD(br BReader, offset uint32, tsg TagSpace, fts FieldTypeSet) (out IFD, err error) {
	if br == nil {
		return nil, errors.New("No BReader supplied.")
	}
	if fts == nil {
		fts = DefaultFieldTypes
	}
	if tsg == nil {
		tsg = DefaultTagSpace
	}
	ifd := new(imageFileDirectory)
	br.Seek(int64(offset), 0)
	if err = br.BRead(&ifd.numEntries); err != nil {
		return
	}
	for i := uint16(0); i < ifd.numEntries; i++ {
		var f Field
		if f, err = ParseField(br, tsg, fts); err != nil {
			return
		}
		ifd.fields = append(ifd.fields, f)
	}
	if err = br.BRead(&ifd.nextOffset); err != nil {
		return
	}
	return ifd, nil
}

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

func ParseIFD8(br BReader, offset uint64, tsg TagSpace, fts FieldTypeSet) (out IFD8, err error) {
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
