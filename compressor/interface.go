package compressor

// Provider defines the basic interface for a data compression element
type Provider interface {
	// Compress will return the compacted version of src data
	Compress(src []byte) []byte

	// Decompress will return the uncompacted version of src data, in case of
	// error the original data will be returned
	Decompress(src []byte) []byte
}

// New is a constructor/setup method to easily create a compressor instance
func New(t string) Provider {
	switch t {
	case "snappy":
		return &snappyCompressor{}
	default:
		return &snappyCompressor{}
	}
}
