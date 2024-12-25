package wolf3d

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestWolfenstein3DMapHeader(t *testing.T) {
	mapHeader, err := readMapHeader(wolfenstein3DMapHeaderData)
	assert.NoError(t, err)

	expectedLevelCount := 10

	assert.NotNil(t, mapHeader)

	assert.Equal(t, uint16(0xABCD), mapHeader.magic)

	assert.Equal(t, expectedLevelCount, mapHeader.LevelCount())
	assert.Greater(t, mapHeader.LevelOffset(expectedLevelCount-1), int32(0))
	assert.Equal(t, int32(-1), mapHeader.LevelOffset(expectedLevelCount))
}
