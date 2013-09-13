package tiff

import (
	"fmt"
	"sync"
)

type TagSet interface {
	RegisterTag(t Tag) error
	GetTag(id uint16) Tag
}

type tagSet struct {
	mu   sync.Mutex
	tags map[uint16]Tag
}

func (ts *tagSet) RegisterTag(t Tag) error {
	return nil
}

func (ts *tagSet) GetTag(id uint16) Tag {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	if t, ok := ts.tags[id]; ok {
		return t
	}
	return &tag{
		id:   id,
		name: fmt.Sprintf("UnregisteredTag_%d", id),
	}
}

var defaultTags = &tagSet{
	tags: map[uint16]Tag{
		// Baseline TIFF Tags
		254:   tagNewSubFileType,
		255:   tagSubfileType,
		256:   tagImageWidth,
		257:   tagImageLength,
		258:   tagBitsPerSample,
		259:   tagCompression,
		262:   tagPhotometricInterpretation,
		263:   tagThreshholding,
		264:   tagCellWidth,
		265:   tagCellLength,
		266:   tagFillOrder,
		270:   tagImageDescription,
		271:   tagMake,
		272:   tagModel,
		273:   tagStripOffsets,
		274:   tagOrientation,
		277:   tagSamplesPerPixel,
		278:   tagRowsPerStrip,
		279:   tagStripByteCounts,
		280:   tagMinSampleValue,
		281:   tagMaxSampleValue,
		282:   tagXResolution,
		283:   tagYResolution,
		284:   tagPlanarConfiguration,
		288:   tagFreeOffsets,
		289:   tagFreeByteCounts,
		290:   tagGrayResponseUnit,
		291:   tagGrayResponseCurve,
		296:   tagResolutionUnit,
		305:   tagSoftware,
		306:   tagDateTime,
		315:   tagArtist,
		316:   tagHostComputer,
		320:   tagColorMap,
		338:   tagExtraSamples,
		33432: tagCopyright,

		// Extension Tags
		269:   tagDocumentName,
		285:   tagPageName,
		286:   tagXPosition,
		287:   tagYPosition,
		292:   tagT4Options,
		293:   tagT6Options,
		297:   tagPageNumber,
		301:   tagTransferFunction,
		317:   tagPredictor,
		318:   tagWhitePoint,
		319:   tagPrimaryChromaticities,
		321:   tagHalftoneHints,
		322:   tagTileWidth,
		323:   tagTileLength,
		324:   tagTileOffsets,
		325:   tagTileByteCounts,
		326:   tagBadFaxLines,
		327:   tagCleanFaxData,
		328:   tagConsecutiveBadFaxLines,
		330:   tagSubIFDs,
		332:   tagInkSet,
		333:   tagInkNames,
		334:   tagNumberOfInks,
		336:   tagDotRange,
		337:   tagTargetPrinter,
		339:   tagSampleFormat,
		340:   tagSMinSampleValue,
		341:   tagSMaxSampleValue,
		342:   tagTransferRange,
		343:   tagClipPath,
		344:   tagXClipPathUnits,
		345:   tagYClipPathUnits,
		346:   tagIndexed,
		347:   tagJPEGTables,
		351:   tagOPIProxy,
		400:   tagGlobalParametersIFD,
		401:   tagProfileType,
		402:   tagFaxProfile,
		403:   tagCodingMethods,
		404:   tagVersionYear,
		405:   tagModeNumber,
		433:   tagDecode,
		434:   tagDefaultImageColor,
		512:   tagJPEGProc,
		513:   tagJPEGInterchangeFormat,
		514:   tagJPEGInterchangeFormatLength,
		515:   tagJPEGRestartInterval,
		517:   tagJPEGLosslessPredictors,
		518:   tagJPEGPointTransforms,
		519:   tagJPEGQTables,
		520:   tagJPEGDCTables,
		521:   tagJPEGACTables,
		529:   tagYCbCrCoefficients,
		530:   tagYCbCrSubSampling,
		531:   tagYCbCrPositioning,
		532:   tagReferenceBlackWhite,
		559:   tagStripRowCounts,
		700:   tagXMP,
		32781: tagImageID,
		34732: tagImageLayer,
	},
}

var DefaultTags TagSet = defaultTags

func RegisterTag(t Tag) error {
	return DefaultTags.RegisterTag(t)
}
