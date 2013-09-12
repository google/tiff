package tiff

// These constants represents the first 4 bytes of the file for each kind of
// TIFF along with each byte ordering.  This is mostly useful for registration
// with the "image" package from the Go standard library.
const (
	hdrTIFFBigEndian    = "MM\x00\x2A"
	hdrTIFFLitEndian    = "II\x2A\x00"
	hdrBigTIFFBigEndian = "MM\x00\x2B"
	hdrBigTIFFLitEndian = "II\x2B\x00"
)

// These constants represent the byte order options present at the beginning of
// a TIFF file.
const (
	BigEndian uint16 = 'M'<<8 | 'M'
	LitEndian uint16 = 'I' | 'I'<<8
)

// These constants represent the TIFF file type identifiers.  At present,
// there are values for a TIFF and a BigTIFF.
const (
	TypeTIFF    uint16 = 0x2A
	TypeBigTIFF uint16 = 0x2B
)
