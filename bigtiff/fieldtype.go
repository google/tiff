package bigtiff

import "github.com/jonathanpittman/tiff"

/* Field type definitions
16 = LONG8	64-bit unsigned integer.
17 = SLONG8	64-bit signed integer.
18 = IFD8	64-bit unsigned integer offset value
*/

var (
	FTLong8  = tiff.NewFieldType(16, "LONG8", 8, false, nil)
	FTSLong8 = tiff.NewFieldType(17, "SLONG8", 8, true, nil)
	FTIFD8   = tiff.NewFieldType(18, "IFD8", 8, false, nil)
)

func init() {
	tiff.RegisterFieldType(FTLong8)
	tiff.RegisterFieldType(FTSLong8)
	tiff.RegisterFieldType(FTIFD8)
}
