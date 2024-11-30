package wolf3d

import (
	"fmt"
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

func TestPrintWolfenstein3DLevelMaps(t *testing.T) {
	t.Skip("Skipping test that print all maps")

	levelMaps, err := Wolfenstein3DMap()
	assert.NoError(t, err)

	lookup := []byte(
		" abcdefghijklmnopqrstuvwxyz1234567890ABCDEFGIJKLMN" +
			"OPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz1234567890AB" +
			"CDEFGIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz1" +
			"234567890ABCDEFGIJKLMNOPQRSTUVWXYZabcdefghijklmnop" +
			"qrstuvwxyz1234567890ABCDEFGIJKLMNOPQRSTUVWXYZabcde" +
			"fghijklmnopqrstuvwxyz1234567890ABCDEFGIJKLMNOPQRST" +
			"UVWXYZabcdefghijklmnopqrstuvwxyz1",
	)

	for _, levelMap := range levelMaps {
		for plane := 0; plane <= 2; plane++ {
			if levelMap.planeExist(plane) {
				var minValue, maxValue = 65535, -65535

				fmt.Printf("Level: \"%s\", plane %d, size: %dx%d\n", levelMap.Name, plane, levelMap.Width, levelMap.Height)

				fmt.Print("    ")
				for x := 0; x < levelMap.Width; x++ {
					fmt.Printf("%d", x%10)
				}
				fmt.Println()

				for y := 0; y < levelMap.Height; y++ {
					fmt.Printf("%3d ", y)
					for x := 0; x < levelMap.Width; x++ {
						value := levelMap.Value(plane, x, y)

						minValue = min(minValue, value)
						maxValue = max(maxValue, value)

						fmt.Print(string(lookup[value]))
					}
					fmt.Println()
				}

				fmt.Printf("Values are in range: %d .. %d\n", minValue, maxValue)
			}
		}
	}
}
