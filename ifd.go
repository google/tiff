package tiff

type IFD struct {
	NumEntries uint16
	Fields     []Field
	NextOffset uint32

	// The image data is not a part of the IFD.  It is always referenced by
	// tags in the IFD.  However, each IFD describes at least one image.
	// Here we provide a place to store that data regardless of how it was
	// originally stored in the file.
	ImageData []byte

	// SubIFDs represents any IFDs that are considered SubIFDs of a top
	// level IFD.  Not all implementations need or use this.  This is here
	// for convenience since SubIFDs do not directly belong to the parent
	// TIFF.  In normal processing, if you only look for the next offset,
	// SubIFDs would be missed.
	SubIFDs  []*IFD
	PrivIFDs []*IFD
}

func (ifd *IFD) processSubIFDs(br *bReader) error {
	return nil
}

func (ifd *IFD) processPrivIFDs(br *bReader) error {
	return nil
}

func (ifd *IFD) processImageData(br *bReader) error {
	return nil
}

func parseIFD(br *bReader, offset uint32) (out *IFD, err error) {
	ifd := new(IFD)
	br.Seek(int64(offset), 0)
	if err = br.Read(&ifd.NumEntries); err != nil {
		return
	}
	for i := uint16(0); i < ifd.NumEntries; i++ {
		var f Field
		if f, err = parseField(br); err != nil {
			return
		}
		ifd.Fields = append(ifd.Fields, f)
	}
	if err = br.Read(&ifd.NextOffset); err != nil {
		return
	}
	// TODO: Look for the image data and process it.
	// TODO: Look for SubIFDs and process them.
	return ifd, nil
}

type IFD8 struct {
	NumEntries uint64
	Fields     []Field8
	NextOffset uint64

	// The image data is not a part of the IFD.  It is always referenced by
	// tags in the IFD.  However, each IFD describes at least one image.
	// Here we provide a place to store that data regardless of how it was
	// originally stored in the file.
	ImageData []byte

	// SubIFDs represents any IFDs that are considered SubIFDs of a top
	// level IFD.  Not all implementations need or use this.  This is here
	// for convenience since SubIFDs do not belong to the parent TIFF.
	SubIFDs  []*IFD8
	PrivIFDs []*IFD8
}

func (ifd8 *IFD8) processImageData(br *bReader) error {
	return nil
}

func parseIFD8(br *bReader, offset uint64) (out *IFD8, err error) {
	ifd := new(IFD8)
	br.Seek(int64(offset), 0) // TODO: This is wrong.  Use big.Int?
	if err = br.Read(&ifd.NumEntries); err != nil {
		return
	}
	for i := uint64(0); i < ifd.NumEntries; i++ {
		var f Field8
		if f, err = parseField8(br); err != nil {
			return
		}
		ifd.Fields = append(ifd.Fields, f)
	}
	if err = br.Read(&ifd.NextOffset); err != nil {
		return
	}
	// TODO: Look for the image data and process it.
	// TODO: Look for SubIFDs and process them.
	return ifd, nil
}
