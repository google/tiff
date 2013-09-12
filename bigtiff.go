package tiff

type BigTIFF struct {
	ByteOrder   uint16 // "MM" or "II"
	Type        uint16 // Must be 43 (0x2B)
	OffsetSize  uint16 // Size in bytes used for offset values
	Constant    uint16 // Must be 0
	FirstOffset uint64 // Offset location for IFD 0
	IFDs        []*IFD8
}
