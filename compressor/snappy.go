package compressor

import "github.com/golang/snappy"

// Snappy-based compressor instance; optimized for speed instead of size
type snappyCompressor struct{}

func (c *snappyCompressor) Compress(src []byte) []byte {
	return snappy.Encode(nil, src)
}

func (c *snappyCompressor) Decompress(src []byte) []byte {
	r, err := snappy.Decode(nil, src)
	if err != nil {
		return src
	}
	return r
}
