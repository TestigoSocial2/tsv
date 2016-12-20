package compressor

import (
	"bytes"
	"testing"
)

func TestSnappyCompress(t *testing.T) {
	data := `Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.`
	c := &snappyCompressor{}
	comp := c.Compress([]byte(data))
	if len(comp) >= len([]byte(data)) {
		t.Error("Unexpected result")
	}
}

func TestSnappyDecompress(t *testing.T) {
	c := &snappyCompressor{}
	comp := c.Compress([]byte("sample data here"))
	res := c.Decompress(comp)
	if string(res) != "sample data here" {
		t.Error("Unexpected result")
	}

	// Invalid decompression
	res2 := c.Decompress([]byte{})
	if !bytes.Equal(res2, []byte{}) {
		t.Error("Unexpected result")
	}
}
