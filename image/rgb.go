package image

/* Baseline RGB

Differences from Palette Color Images

Tag 258 (BitsPerSample)
	BitsPerSample = 8,8,8
	Note: Each component is 8 bits deep in a Baseline TIFF RGB image.

Tag 262 (PhotometricInterpretation)
	PhotometricInterpretation = 2 (RGB)

There is no ColorMap.

Tag 277 (SamplesPerPixel)
	Note: The number of components per pixel. This number is 3 for RGB images, unless extra samples are present. See the ExtraSamples field for further information.

Required Fields
		<Baseline Grayscale>
		ImageWidth
		ImageLength
		BitsPerSample
		Compression
		PhotometricInterpretation
		StripOffsets
	SamplesPerPixel
		RowsPerStrip
		StripByteCounts
		XResolution
		YResolution
		ResolutionUnit

*/

type FullColorRGB struct {
	Grayscale
	SamplesPerPixel []uint16 ``
}
