package wolf3d

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCarmackDecode(t *testing.T) {
	// Test special case for words with high byte values 0xA7 and 0xA8
	// https://moddingwiki.shikadi.net/wiki/Carmack_compression

	type testCase struct {
		name         string
		compressed   []byte
		decompressed []byte
	}

	tests := []testCase{
		{name: "special case: 0xA7 as high byte values", compressed: []byte{0x00, 0xA7, 0x12, 0xEE, 0xFF, 0x00, 0xA8, 0x34, 0xCC, 0xDD}, decompressed: []byte{0x12, 0xA7, 0xEE, 0xFF, 0x34, 0xA8, 0xCC, 0xDD}},
		{name: "special case: 0xA8 as high byte values", compressed: []byte{0x00, 0xA7, 0x12, 0xEE, 0xFF, 0x00, 0xA7, 0x34, 0xCC, 0xDD}, decompressed: []byte{0x12, 0xA7, 0xEE, 0xFF, 0x34, 0xA7, 0xCC, 0xDD}},
		{name: "near pointer before EOF", compressed: []byte{0x78, 0x56, 0x34, 0x12, 0x02, 0xA7, 0x02, 0x00, 0x01}, decompressed: []byte{0x78, 0x56, 0x34, 0x12, 0x78, 0x56, 0x34, 0x12, 0x00, 0x01}},
		{name: "near pointer at EOF", compressed: []byte{0x78, 0x56, 0x34, 0x12, 0x02, 0xA7, 0x02}, decompressed: []byte{0x78, 0x56, 0x34, 0x12, 0x78, 0x56, 0x34, 0x12}},
		{name: "far pointer before EOF", compressed: []byte{0x78, 0x56, 0x34, 0x12, 0x01, 0xA8, 0x01, 0x00, 0x01, 0x02}, decompressed: []byte{0x78, 0x56, 0x34, 0x12, 0x34, 0x12, 0x01, 0x02}},
		{name: "far pointer at EOF", compressed: []byte{0x78, 0x56, 0x34, 0x12, 0x01, 0xA8, 0x01, 0x00}, decompressed: []byte{0x78, 0x56, 0x34, 0x12, 0x34, 0x12}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			decompressed := CarmackDecode(tt.compressed)

			assert.Equal(t, len(tt.decompressed), len(decompressed))
			assert.Equal(t, tt.decompressed, decompressed)
		})
	}
}

func TestCarmackDecodeWolf3dLevel1Plane0(t *testing.T) {
	compressedData, err := readFile("testdata/wolfenstein_level1_plane0_RLEW_Carmack_compressed.bin")
	assert.NoError(t, err)

	size, decompressedData := CarmackDecodeWithLengthPrefix(compressedData)

	expectedData, err := readFile("testdata/wolfenstein_level1_plane0_RLEW_compressed.bin")
	assert.NoError(t, err)
	assert.Equal(t, size, len(decompressedData))
	assert.Equal(t, expectedData, decompressedData)
}
