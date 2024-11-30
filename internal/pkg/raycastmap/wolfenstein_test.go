package raycastmap

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"maze/internal/pkg/wolf3d"
	"slices"
	"testing"
)

func TestPrintWolfensteinMap(t *testing.T) {
	levelMaps, err := wolf3d.Wolfenstein3DMap()
	assert.NoError(t, err)

	level := levelMaps[0]
	mazeWallPlane := 0

	fmt.Println(level.Name)
	for y := 0; y < level.Height; y++ {
		for x := 0; x < level.Width; x++ {
			value := level.Value(mazeWallPlane, x, y)

			if slices.Contains([]int{0x6A, 0x6C, 0x6D, 0x6E, 0x6F, 0x70, 0x71, 0x72, 0x75, 0x76, 0x8B, 0x8C, 0x8D, 0x8E, 0x8F}, value) {
				fmt.Print("  ")
				value = 0x00
			} else if slices.Contains([]int{0x5A, 0x5B}, value) {
				fmt.Print("++")
				value = 0x00
			} else {
				fmt.Printf("%02X", value)
			}
		}
		fmt.Println()
	}
}
