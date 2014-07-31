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
		return nil, errors.New("tiff: no BReader supplied")
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
