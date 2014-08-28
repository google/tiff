package bigtiff

const (
	MagicBigEndian        = "MM\x00\x2B"
	MagicLitEndian        = "II\x2B\x00"
	Version        uint16 = 0x2B
	VersionName    string = "BigTIFF"
)
