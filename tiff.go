package tiff

type TIFF struct {
	ByteOrder   uint16 // "MM" or "II"
	Type        uint16 // Must be 42 (0x2A)
	FirstOffset uint32 // Offset location for IFD 0
	IFDs        []*IFD
}
