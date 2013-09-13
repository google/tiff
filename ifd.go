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
	// for convenience since SubIFDs do not belong to the parent TIFF.
	SubIFDs  []*IFD
	PrivIFDs []*IFD
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
