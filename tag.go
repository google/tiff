package tiff

type Tag interface {
	Id() uint16
	Name() string
	ValidFieldTypes() []FieldType
}

type tag struct {
	id       uint16
	name     string
	validFTs []FieldType
}

func (t *tag) Id() uint16 {
	return t.id
}

func (t *tag) Name() string {
	return t.name
}

func (t *tag) ValidFieldTypes() []FieldType {
	return t.validFTs
}

var (
	// Baseline TIFF Tags
	tagNewSubFileType            = &tag{id: 254, name: "NewSubfileType", validFTs: []FieldType{fTLong}}
	tagSubfileType               = &tag{id: 255, name: "SubfileType", validFTs: []FieldType{}}
	tagImageWidth                = &tag{id: 256, name: "ImageWidth", validFTs: []FieldType{}}
	tagImageLength               = &tag{id: 257, name: "ImageLength", validFTs: []FieldType{}}
	tagBitsPerSample             = &tag{id: 258, name: "BitsPerSample", validFTs: []FieldType{}}
	tagCompression               = &tag{id: 259, name: "Compression", validFTs: []FieldType{}}
	tagPhotometricInterpretation = &tag{id: 262, name: "PhotometricInterpretation", validFTs: []FieldType{}}
	tagThreshholding             = &tag{id: 263, name: "Threshholding", validFTs: []FieldType{}}
	tagCellWidth                 = &tag{id: 264, name: "CellWidth", validFTs: []FieldType{}}
	tagCellLength                = &tag{id: 265, name: "CellLength", validFTs: []FieldType{}}
	tagFillOrder                 = &tag{id: 266, name: "FillOrder", validFTs: []FieldType{}}
	tagImageDescription          = &tag{id: 270, name: "ImageDescription", validFTs: []FieldType{}}
	tagMake                      = &tag{id: 271, name: "Make", validFTs: []FieldType{}}
	tagModel                     = &tag{id: 272, name: "Model", validFTs: []FieldType{}}
	tagStripOffsets              = &tag{id: 273, name: "StripOffsets", validFTs: []FieldType{}}
	tagOrientation               = &tag{id: 274, name: "Orientation", validFTs: []FieldType{}}
	tagSamplesPerPixel           = &tag{id: 277, name: "SamplesPerPixel", validFTs: []FieldType{}}
	tagRowsPerStrip              = &tag{id: 278, name: "RowsPerStrip", validFTs: []FieldType{}}
	tagStripByteCounts           = &tag{id: 279, name: "StripByteCounts", validFTs: []FieldType{}}
	tagMinSampleValue            = &tag{id: 280, name: "MinSampleValue", validFTs: []FieldType{}}
	tagMaxSampleValue            = &tag{id: 281, name: "MaxSampleValue", validFTs: []FieldType{}}
	tagXResolution               = &tag{id: 282, name: "XResolution", validFTs: []FieldType{}}
	tagYResolution               = &tag{id: 283, name: "YResolution", validFTs: []FieldType{}}
	tagPlanarConfiguration       = &tag{id: 284, name: "PlanarConfiguration", validFTs: []FieldType{}}
	tagFreeOffsets               = &tag{id: 288, name: "FreeOffsets", validFTs: []FieldType{}}
	tagFreeByteCounts            = &tag{id: 289, name: "FreeByteCounts", validFTs: []FieldType{}}
	tagGrayResponseUnit          = &tag{id: 290, name: "GrayResponseUnit", validFTs: []FieldType{}}
	tagGrayResponseCurve         = &tag{id: 291, name: "GrayResponseCurve", validFTs: []FieldType{}}
	tagResolutionUnit            = &tag{id: 296, name: "ResolutionUnit", validFTs: []FieldType{}}
	tagSoftware                  = &tag{id: 305, name: "Software", validFTs: []FieldType{}}
	tagDateTime                  = &tag{id: 306, name: "DateTime", validFTs: []FieldType{}}
	tagArtist                    = &tag{id: 315, name: "Artist", validFTs: []FieldType{}}
	tagHostComputer              = &tag{id: 316, name: "HostComputer", validFTs: []FieldType{}}
	tagColorMap                  = &tag{id: 320, name: "ColorMap", validFTs: []FieldType{}}
	tagExtraSamples              = &tag{id: 338, name: "ExtraSamples", validFTs: []FieldType{}}
	tagCopyright                 = &tag{id: 33432, name: "Copyright", validFTs: []FieldType{}}

	// Extension Tags
	tagDocumentName                = &tag{id: 269, name: "DocumentName", validFTs: []FieldType{}}
	tagPageName                    = &tag{id: 285, name: "PageName", validFTs: []FieldType{}}
	tagXPosition                   = &tag{id: 286, name: "XPosition", validFTs: []FieldType{}}
	tagYPosition                   = &tag{id: 287, name: "YPosition", validFTs: []FieldType{}}
	tagT4Options                   = &tag{id: 292, name: "T4Options", validFTs: []FieldType{}}
	tagT6Options                   = &tag{id: 293, name: "T6Options", validFTs: []FieldType{}}
	tagPageNumber                  = &tag{id: 297, name: "PageNumber", validFTs: []FieldType{}}
	tagTransferFunction            = &tag{id: 301, name: "TransferFunction", validFTs: []FieldType{}}
	tagPredictor                   = &tag{id: 317, name: "Predictor", validFTs: []FieldType{}}
	tagWhitePoint                  = &tag{id: 318, name: "WhitePoint", validFTs: []FieldType{}}
	tagPrimaryChromaticities       = &tag{id: 319, name: "PrimaryChromaticities", validFTs: []FieldType{}}
	tagHalftoneHints               = &tag{id: 321, name: "HalftoneHints", validFTs: []FieldType{}}
	tagTileWidth                   = &tag{id: 322, name: "TileWidth", validFTs: []FieldType{}}
	tagTileLength                  = &tag{id: 323, name: "TileLength", validFTs: []FieldType{}}
	tagTileOffsets                 = &tag{id: 324, name: "TileOffsets", validFTs: []FieldType{}}
	tagTileByteCounts              = &tag{id: 325, name: "TileByteCounts", validFTs: []FieldType{}}
	tagBadFaxLines                 = &tag{id: 326, name: "BadFaxLines", validFTs: []FieldType{}}
	tagCleanFaxData                = &tag{id: 327, name: "CleanFaxData", validFTs: []FieldType{}}
	tagConsecutiveBadFaxLines      = &tag{id: 328, name: "ConsecutiveBadFaxLines", validFTs: []FieldType{}}
	tagSubIFDs                     = &tag{id: 330, name: "SubIFDs", validFTs: []FieldType{}}
	tagInkSet                      = &tag{id: 332, name: "InkSet", validFTs: []FieldType{}}
	tagInkNames                    = &tag{id: 333, name: "InkNames", validFTs: []FieldType{}}
	tagNumberOfInks                = &tag{id: 334, name: "NumberOfInks", validFTs: []FieldType{}}
	tagDotRange                    = &tag{id: 336, name: "DotRange", validFTs: []FieldType{}}
	tagTargetPrinter               = &tag{id: 337, name: "TargetPrinter", validFTs: []FieldType{}}
	tagSampleFormat                = &tag{id: 339, name: "SampleFormat", validFTs: []FieldType{}}
	tagSMinSampleValue             = &tag{id: 340, name: "SMinSampleValue", validFTs: []FieldType{}}
	tagSMaxSampleValue             = &tag{id: 341, name: "SMaxSampleValue", validFTs: []FieldType{}}
	tagTransferRange               = &tag{id: 342, name: "TransferRange", validFTs: []FieldType{}}
	tagClipPath                    = &tag{id: 343, name: "ClipPath", validFTs: []FieldType{}}
	tagXClipPathUnits              = &tag{id: 344, name: "XClipPathUnits", validFTs: []FieldType{}}
	tagYClipPathUnits              = &tag{id: 345, name: "YClipPathUnits", validFTs: []FieldType{}}
	tagIndexed                     = &tag{id: 346, name: "Indexed", validFTs: []FieldType{}}
	tagJPEGTables                  = &tag{id: 347, name: "JPEGTables", validFTs: []FieldType{}}
	tagOPIProxy                    = &tag{id: 351, name: "OPIProxy", validFTs: []FieldType{}}
	tagGlobalParametersIFD         = &tag{id: 400, name: "GlobalParametersIFD", validFTs: []FieldType{}}
	tagProfileType                 = &tag{id: 401, name: "ProfileType", validFTs: []FieldType{}}
	tagFaxProfile                  = &tag{id: 402, name: "FaxProfile", validFTs: []FieldType{}}
	tagCodingMethods               = &tag{id: 403, name: "CodingMethods", validFTs: []FieldType{}}
	tagVersionYear                 = &tag{id: 404, name: "VersionYear", validFTs: []FieldType{}}
	tagModeNumber                  = &tag{id: 405, name: "ModeNumber", validFTs: []FieldType{}}
	tagDecode                      = &tag{id: 433, name: "Decode", validFTs: []FieldType{}}
	tagDefaultImageColor           = &tag{id: 434, name: "DefaultImageColor", validFTs: []FieldType{}}
	tagJPEGProc                    = &tag{id: 512, name: "JPEGProc", validFTs: []FieldType{}}
	tagJPEGInterchangeFormat       = &tag{id: 513, name: "JPEGInterchangeFormat", validFTs: []FieldType{}}
	tagJPEGInterchangeFormatLength = &tag{id: 514, name: "JPEGInterchangeFormatLength", validFTs: []FieldType{}}
	tagJPEGRestartInterval         = &tag{id: 515, name: "JPEGRestartInterval", validFTs: []FieldType{}}
	tagJPEGLosslessPredictors      = &tag{id: 517, name: "JPEGLosslessPredictors", validFTs: []FieldType{}}
	tagJPEGPointTransforms         = &tag{id: 518, name: "JPEGPointTransforms", validFTs: []FieldType{}}
	tagJPEGQTables                 = &tag{id: 519, name: "JPEGQTables", validFTs: []FieldType{}}
	tagJPEGDCTables                = &tag{id: 520, name: "JPEGDCTables", validFTs: []FieldType{}}
	tagJPEGACTables                = &tag{id: 521, name: "JPEGACTables", validFTs: []FieldType{}}
	tagYCbCrCoefficients           = &tag{id: 529, name: "YCbCrCoefficients", validFTs: []FieldType{}}
	tagYCbCrSubSampling            = &tag{id: 530, name: "YCbCrSubSampling", validFTs: []FieldType{}}
	tagYCbCrPositioning            = &tag{id: 531, name: "YCbCrPositioning", validFTs: []FieldType{}}
	tagReferenceBlackWhite         = &tag{id: 532, name: "ReferenceBlackWhite", validFTs: []FieldType{}}
	tagStripRowCounts              = &tag{id: 559, name: "StripRowCounts", validFTs: []FieldType{}}
	tagXMP                         = &tag{id: 700, name: "XMP", validFTs: []FieldType{}}
	tagImageID                     = &tag{id: 32781, name: "ImageID", validFTs: []FieldType{}}
	tagImageLayer                  = &tag{id: 34732, name: "ImageLayer", validFTs: []FieldType{}}

	// Private Tags
)
