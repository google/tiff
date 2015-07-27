// Package modi provides tiff extensions for working with Microsoft Office
// Document Imaging based tiff files.
package modi // import "jonathanpittman.com/tiff/modi"

import "jonathanpittman.com/tiff"

var modiTags = tiff.NewTagSet("MODI", 32768, 65535)

func init() {
	modiTags.Register(tiff.NewTag(37679, "MODIText", nil))
	modiTags.Register(tiff.NewTag(37680, "MODIOLEPropertySetStorage", nil))
	modiTags.Register(tiff.NewTag(37681, "MODIPositioning", nil))

	modiTags.Lock()

	tiff.DefaultTagSpace.RegisterTagSet(modiTags)
}
