package tiff

import (
	"io"
)

// ReadAtReadSeeker is the interface that wraps the Read, ReadAt, and Seek
// methods.  Typical use cases would satisfy this with a bytes.Reader (in
// memory) or an os.File (on disk).  For truly large files, such as BigTIFF, a
// user may want to create a custom solution that combines both in memory and on
// disk solutions for accessing the contents.
type ReadAtReadSeeker interface {
	io.ReadSeeker
	io.ReaderAt
}
