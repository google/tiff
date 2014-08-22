package image

/* Baseline Grayscale

Differences from Bilevel Images
	Change:
		Tag 259 (Compression)
			Note: Must be 1 or 32773 only

	Add:
		Tag 258 (BitsPerSample)
			4 = 16 shades of gray
			8 = 256 shades of gray
			Note: The number of bits per component.
			Note: Allowable values for Baseline TIFF grayscale
				images are 4 and 8.

Required Fields
		<Baseline Bilevel>
		ImageWidth
		ImageLength
	BitsPerSample
		Compression
		PhotometricInterpretation
		StripOffsets
		RowsPerStrip
		StripByteCounts
		XResolution
		YResolution
		ResolutionUnit
*/

type Grayscale struct {
	Bilevel
	BitsPerSample []uint16 `tifftag:"id=258"`
}
