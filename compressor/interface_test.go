package compressor

import "testing"

func TestNew(t *testing.T) {
	if _, ok := New("snappy").(*snappyCompressor); !ok {
		t.Error("Unexpected compressor type")
	}

	if _, ok := New("invalid").(*snappyCompressor); !ok {
		t.Error("Unexpected compressor type")
	}
}
