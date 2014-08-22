package image

/* Baseline Palette-Color

Differences from Grayscale Images
	Tag 262 (PhotometricInterpretation)
		PhotometricInterpretation = 3 (Palette Color)

	Tag 320 (ColorMap)
		Note: 3 * (2**BitsPerSample)
		Note: This field defines a Red-Green-Blue color map (often
			called a lookup table) for palette color images. In a
			palette-color image, a pixel value is used to index into
			an RGB-lookup table. For example, a palette-color pixel
			having a value of 0 would be displayed according to the
			0th Red, Green, Blue triplet.
			In a TIFF ColorMap, all the Red values come first,
			followed by the Green values, then the Blue values. In
			the ColorMap, black is represented by 0,0,0 and white is
			represented by 65535, 65535, 65535.

Required Fields
		<Baseline Grayscale>
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
	ColorMap
*/

type PaletteColor struct {
	Grayscale
	ColorMap []uint16 `tifftag:"id=320"`
}
