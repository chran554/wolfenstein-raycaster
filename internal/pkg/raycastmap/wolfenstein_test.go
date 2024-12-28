package raycastmap

import (
	"fmt"
	"github.com/anthonynsimon/bild/blend"
	"github.com/stretchr/testify/assert"
	"image"
	"image/color"
	"image/png"
	"maze/internal/pkg/wolf3d"
	"os"
	"slices"
	"strconv"
	"strings"
	"testing"
)

func TestPrintWolfensteinMapValues(t *testing.T) {
	t.Skip("Skipping test that print map to console")

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
	t.Skip("Skipping test that print map to console")

	levelMaps, err := wolf3d.Wolfenstein3DMap()
	assert.NoError(t, err)

	level := 0                // Level to render
	const showSpecials = true // Show special locations, items, foes, start point et.c.

	levelMap := levelMaps[level]

	const wallDataPlane = 0
	const itemsDataPlane = 1

	startPoint := Cell{Structure: StructureNone}

	fmt.Println("Level name:", levelMap.Name)

	fmt.Print("   ")
	for x := 0; x < levelMap.Width; x++ {
		fmt.Printf("%2d", x%10)
	}
	fmt.Println()

	for y := 0; y < levelMap.Height; y++ {
		fmt.Printf("%2d  ", y)

		for x := 0; x < levelMap.Width; x++ {
			cellValue := levelMap.Value(wallDataPlane, x, y)
			specialValue := levelMap.Value(itemsDataPlane, x, y)

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

func TestPaintWolfensteinMap(t *testing.T) {
	t.Skip("Skipping test that creates image file in project directory")

	level := 0

	cellWidth := 64

	levelMap, err := NewWolfensteinMap(level)
	assert.NoError(t, err)

	mapImage := image.NewRGBA(image.Rect(level, level, levelMap.Width()*(cellWidth+1)+1, levelMap.Height()*(cellWidth+1)+1))

	UnknownTexture := NewTextureFromFile("overlay/question-mark.png")

	startCell := Cell{Structure: StructureNone}

	for y := level; y < levelMap.Height(); y++ {
		for x := level; x < levelMap.Width(); x++ {
			structure := levelMap.StructureAt(x, y)
			special := levelMap.SpecialAt(x, y)

			if special == SpecialStartPointFacingNorth || special == SpecialStartPointFacingSouth || special == SpecialStartPointFacingEast || special == SpecialStartPointFacingWest {
				startCell = Cell{X: x, Y: y, Structure: special}
			}

			if structure == StructureExitDoor && structure.Texture2 != nil {
				// Exit door icon uses texture 2 (exit room sides are texture 1)
				exitTexture := structure.Texture.img // Elevator handle bars walls

				noWestWall := levelMap.StructureAt(x-1, y) == StructureNone
				noEastWall := levelMap.StructureAt(x+1, y) == StructureNone
				if noWestWall || noEastWall {
					exitTexture = structure.Texture2.img // Elevator control panel wall
				}

				renderMapIcon(mapImage, x, levelMap.Height()-1-y, cellWidth, exitTexture)
			} else if structure != StructureNone && structure.Texture != nil {
				renderMapIcon(mapImage, x, levelMap.Height()-1-y, cellWidth, structure.Texture.img)
			}

			if special != SpecialNone && special.Texture != nil {
				renderMapIcon(mapImage, x, levelMap.Height()-1-y, cellWidth, special.Texture.img)
			} else if special != SpecialNone {
				renderMapIcon(mapImage, x, levelMap.Height()-1-y, cellWidth, UnknownTexture.img)
			}
		}
	}

	renderMapIcon(mapImage, startCell.X, levelMap.Height()-1-startCell.Y, cellWidth, SpecialBlueOrb.Texture.img)

	renderCellBorders(mapImage, levelMap.Width(), levelMap.Height(), cellWidth)

	// Save map image
	filename := fmt.Sprintf("wolfenstein3d-map-level%d.png", level)
	imageFile, err := os.Create(filename)
	assert.NoError(t, err)
	err = png.Encode(imageFile, mapImage)
	assert.NoError(t, err)
	err = imageFile.Close()
	assert.NoError(t, err)

	t.Logf("Wrote Wolfenstein 3D map (level %d) to file: %s\n", level, filename)
}

func renderCellBorders(mapImage *image.RGBA, width int, height int, cellWidth int) {
	c := color.NRGBA{R: 107, G: 107, B: 107, A: 255} // Same color as the floor in the maze

	for cellY := 0; cellY < height; cellY++ {
		for cellX := 0; cellX < width; cellX++ {
			for pixel := 0; pixel < cellWidth+1; pixel++ {
				mapImage.Set(cellX*(cellWidth+1)+pixel, cellY*(cellWidth+1), c)
				mapImage.Set(cellX*(cellWidth+1), cellY*(cellWidth+1)+pixel, c)

				if cellX == width-1 {
					mapImage.Set((cellX+1)*(cellWidth+1), cellY*(cellWidth+1)+pixel, c)
				}
				if cellY == height-1 {
					mapImage.Set(cellX*(cellWidth+1)+pixel, (cellY+1)*(cellWidth+1), c)
				}
			}
		}
	}
}

func renderMapIcon(mapImage *image.RGBA, cellX int, cellY int, cellSize int, texture image.Image) {
	textureWidth := texture.Bounds().Dx()
	textureHeight := texture.Bounds().Dy()

	startX := cellX*(cellSize+1) + 1 + (cellSize-textureWidth)/2
	startY := cellY*(cellSize+1) + 1 + (cellSize-textureHeight)/2

	subMapImage := mapImage.SubImage(image.Rect(startX, startY, startX+textureWidth, startY+textureHeight))
	resultImage := blend.Normal(subMapImage, texture)

	for y := 0; y < textureHeight; y++ {
		for x := 0; x < textureWidth; x++ {
			mapImage.Set(startX+x, startY+y, resultImage.At(x, y))
		}
	}
}
