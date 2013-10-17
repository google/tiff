package tiff

import (
	"fmt"
	"sync"
)

type TagSet interface {
	GetTag(id uint16) (Tag, bool)
	Name() string
}

type tagSet struct {
	name string
	tags map[uint16]Tag
}

func (ts *tagSet) GetTag(id uint16) (Tag, bool) {
	t, ok := ts.tags[id]
	return t, ok
}

func (ts *tagSet) Name() string {
	return ts.name
}

var baselineTags = &tagSet{
	name: "Baseline",
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
	},
}

var extendedTags = &tagSet{
	name: "Extended",
	tags: map[uint16]Tag{
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

type TagSetGroup interface {
	Name() string
	RegisterTagSet(ts TagSet) int
	GetTagSet(id int) (TagSet, error)
	GetTag(id uint16) Tag
}

type tagSetGroup struct {
	mu   sync.Mutex
	name string
	ts   []TagSet
}

func (tsg *tagSetGroup) Name() string {
	return tsg.name
}

func (tsg *tagSetGroup) RegisterTagSet(ts TagSet) int {
	tsg.mu.Lock()
	defer tsg.mu.Unlock()
	tsg.ts = append(tsg.ts, ts)
	return len(tsg.ts) - 1
}

func (tsg *tagSetGroup) GetTagSet(id int) (TagSet, error) {
	tsg.mu.Lock()
	defer tsg.mu.Unlock()
	if id >= len(tsg.ts) {
		return nil, fmt.Errorf("tiff: no TagSet with an id of %d", id)
	}
	return tsg.ts[id], nil
}

func (tsg *tagSetGroup) GetTag(id uint16) Tag {
	tsg.mu.Lock()
	defer tsg.mu.Unlock()
	for _, ts := range tsg.ts {
		if t, ok := ts.GetTag(id); ok {
			return t
		}
	}
	return &tag{id: id}
}

var defTagSetGroup = &tagSetGroup{
	name: "DefaultTagSetGroup",
	ts:   []TagSet{baselineTags, extendedTags},
}

var DefaultTagSetGroup TagSetGroup = defTagSetGroup

func RegisterTagSet(ts TagSet) int {
	return DefaultTagSetGroup.RegisterTagSet(ts)
}
