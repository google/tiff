package exif

import "github.com/jonathanpittman/tiff"

var (
	iopTags     = tiff.NewTagSet("Interoperability", 0, 65535)
	IOPTagSpace = tiff.NewTagSpace("Interoperability")
	iopIFDTag   = tiff.NewTag(40965, "Interoperability IFD", nil)
)

func init() {
	tiff.PrivateTags.Register(iopIFDTag)

	// http://www.exiv2.org/tags.html
	iopTags.Register(tiff.NewTag(1, "InteroperabilityIndex", nil))
	iopTags.Register(tiff.NewTag(2, "InteroperabilityVersion", nil))
	iopTags.Register(tiff.NewTag(4096, "RelatedImageFileFormat", nil))
	iopTags.Register(tiff.NewTag(4097, "RelatedImageWidth", nil))
	iopTags.Register(tiff.NewTag(4098, "RelatedImageLength", nil))

	iopTags.Lock()

	IOPTagSpace.RegisterTagSet(iopTags)
	tiff.RegisterTagSpace(IOPTagSpace)
}
