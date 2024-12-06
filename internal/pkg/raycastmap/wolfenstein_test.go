package raycastmap

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"maze/internal/pkg/wolf3d"
	"slices"
	"strconv"
	"strings"
	"testing"
)

func TestPrintWolfensteinMapValues(t *testing.T) {
	levelMaps, err := wolf3d.Wolfenstein3DMap()
	assert.NoError(t, err)

	walls := make([]string, 256)
	special := make([]string, 256)

	const mazeWallPlane = 0
	const mazeSpecialsPlane = 1

	for levelIndex, levelMap := range levelMaps {
		startPoint := Cell{Structure: StructureNone}

		for y := 0; y < levelMap.Height; y++ {
			for x := 0; x < levelMap.Width; x++ {
				cellValue := levelMap.Value(mazeWallPlane, x, y)
				specialValue := levelMap.Value(mazeSpecialsPlane, x, y)

				// Record location of start point
				if specialValue == 0x14 || specialValue == 0x13 || specialValue == 0x15 {
					startPoint.X = x
					startPoint.Y = y
				}

				levelIndexText := strconv.Itoa(levelIndex)

				if !strings.HasSuffix(walls[cellValue], levelIndexText) {
					walls[cellValue] += " " + levelIndexText
				}

				if !strings.HasSuffix(special[specialValue], levelIndexText) {
					special[specialValue] += " " + levelIndexText
				}

			}
		}
		fmt.Println("Level ", levelIndex, " start point: ", startPoint.X, startPoint.Y)
	}

	fmt.Println()

	fmt.Println("Wall values and levels they are used in:")
	for wallIndex, levels := range walls {
		if levels != "" {
			fmt.Printf("0x%02X   %s\n", wallIndex, levels)
		}
	}

	fmt.Println()

	fmt.Println("Special values and levels they are used in:")
	for specialIndex, levels := range special {
		if levels != "" {
			fmt.Printf("0x%02X   %s\n", specialIndex, levels)
		}
	}

}

func TestPrintWolfensteinMap(t *testing.T) {
	levelMaps, err := wolf3d.Wolfenstein3DMap()
	assert.NoError(t, err)

	level := 3                // Level to render
	const showSpecials = true // Show special locations, items, foes, start point et.c.

	levelMap := levelMaps[level]

	const mazeWallPlane = 0
	const mazeSpecialsPlane = 1

	startPoint := Cell{Structure: StructureNone}

	fmt.Println(levelMap.Name)
	for y := 0; y < levelMap.Height; y++ {
		for x := 0; x < levelMap.Width; x++ {
			cellValue := levelMap.Value(mazeWallPlane, x, y)
			specialValue := levelMap.Value(mazeSpecialsPlane, x, y)

			// Record location of start point
			if specialValue == 0x14 || specialValue == 0x13 || specialValue == 0x15 {
				startPoint.X = x
				startPoint.Y = y
			}

			emptyCellValues := []int{
				0x6A, 0x6C, 0x6D, 0x6E, 0x6F,
				0x70, 0x71, 0x72, 0x73, 0x74, 0x75, 0x76, 0x77, 0x78, 0x79, 0x7A, 0x7B, 0x7C, 0x7D, 0x7E, 0x7F,
				0x82, 0x83, 0x84, 0x85, 0x86, 0x87, 0x88, 0x8A, 0x8B, 0x8C, 0x8D, 0x8E, 0x8F}

			valueText := "  "
			if slices.Contains(emptyCellValues, cellValue) {
				if showSpecials {
					if specialValue != 0x00 {
						valueText = fmt.Sprintf("%02X", specialValue)
					}
				}
			} else if cellValue == 0x5A {
				valueText = fmt.Sprint("| ")
			} else if cellValue == 0x5B {
				valueText = fmt.Sprint("--")
			} else {
				if showSpecials {
					valueText = "##"
				} else {
					valueText = fmt.Sprintf("%02X", cellValue) // The value of the wall
				}
			}

			fmt.Print(valueText)
		}
		fmt.Println()
	}
	fmt.Println("Start point at: ", startPoint.X, startPoint.Y)
}
