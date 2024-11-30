package wolf3d

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRLEWDecodeWolf3dLevel1Plane0(t *testing.T) {
	compressedData, err := readFile("testdata/wolfenstein_level1_plane0_RLEW_compressed.bin")
	assert.NoError(t, err)

	const rleFlag = 0xABCD // Is the "magic" number from the MAPHEAD.WL1 map header file. It is used for RLE expansion.
	size, decompressedData := RLEWDecodeWithLengthPrefixAndRLEFlag(compressedData, rleFlag)

	expectedData, err := readFile("testdata/wolfenstein_level1_plane0_level_data.bin")
	assert.NoError(t, err)
	assert.Equal(t, size, len(decompressedData))
	assert.Equal(t, expectedData, decompressedData)
}
