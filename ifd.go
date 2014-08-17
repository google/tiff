package tiff

import (
	"bytes"
	"errors"
	"fmt"
	"text/tabwriter"
)

// IFD represents the data structure of an IFD in a TIFF File.
type IFD interface {
	NumEntries() uint16
	Fields() []Field
	NextOffset() uint32
	HasField(id uint16) bool
	GetField(id uint16) Field
}

type imageFileDirectory struct {
	numEntries uint16
	fields     []Field
	nextOffset uint32
	fieldMap   map[uint16]Field
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

func (ifd *imageFileDirectory) HasField(id uint16) bool {
	_, ok := ifd.fieldMap[id]
	return ok
}

func (ifd *imageFileDirectory) GetField(id uint16) Field {
	return ifd.fieldMap[id]
}

func (ifd *imageFileDirectory) String() string {
	fmtStr := `
NumEntries: %d
NextOffset: %d
Fields (%d):
%s
`
	w := new(tabwriter.Writer)
	var buf bytes.Buffer

	w.Init(&buf, 5, 0, 1, ' ', 0)
	for _, f := range ifd.Fields() {
		fmt.Fprintf(w, "  %v\n", f)
	}
	w.Flush()

	return fmt.Sprintf(fmtStr, ifd.numEntries, ifd.nextOffset, len(ifd.fields), buf.String())
}

func ParseIFD(br BReader, offset uint32, tsp TagSpace, ftsp FieldTypeSpace) (out IFD, err error) {
	if br == nil {
		return nil, errors.New("tiff: no BReader supplied")
	}
	if ftsp == nil {
		ftsp = DefaultFieldTypeSpace
	}
	if tsp == nil {
		tsp = DefaultTagSpace
	}
	ifd := &imageFileDirectory{
		fieldMap: make(map[uint16]Field, 1),
	}
	br.Seek(int64(offset), 0)
	if err = br.BRead(&ifd.numEntries); err != nil {
		err = fmt.Errorf("tiff: unable to read the number of entries for the IFD at offset %#08x: %v", offset, err)
		return
	}
	for i := uint16(0); i < ifd.numEntries; i++ {
		var f Field
		if f, err = ParseField(br, tsp, ftsp); err != nil {
			return
		}
		ifd.fields = append(ifd.fields, f)
		ifd.fieldMap[f.Tag().ID()] = f
	}
	if err = br.BRead(&ifd.nextOffset); err != nil {
		err = fmt.Errorf("tiff: unable to read the offset for the next ifd: %v", err)
		return
	}
	return ifd, nil
}
