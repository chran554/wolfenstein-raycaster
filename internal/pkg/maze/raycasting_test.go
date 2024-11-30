package maze

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"math"
	"maze/internal/pkg/raycastmap"
	"strings"
	"testing"
)

func TestZeroDivision(t *testing.T) {
	var u = 100.0
	var n = 0.0

	v := u / n

	assert.Equal(t, math.Inf(1), v)
}

func wallValueToStructure(v int) raycastmap.Structure {
	if v == 0 {
		return raycastmap.StructureNone
	}

	return raycastmap.StructureGreyStoneWall1
}

func TestRaycastDistance(t *testing.T) {
	horizontalTestMap := raycastmap.NewSliceMap([][]int{
		{1}, {0}, {0}, {0}, {0}, {0}, {0}, {0}, {0}, {1},
	}, 0.0, 0.0, 0.0, wallValueToStructure)

	horizontalIntersectionInfo1 := RaycastRay(&Vector{4.5, 0.5}, &Vector{-1.0, 0.0}, horizontalTestMap)
	horizontalIntersectionInfo2 := RaycastRay(&Vector{4.5, 0.5}, &Vector{1.0, 0.0}, horizontalTestMap)

	assert.Equal(t, 1.0, horizontalIntersectionInfo1.IntersectionPoint.X)
	assert.Equal(t, 0.5, horizontalIntersectionInfo1.IntersectionPoint.Y)
	assert.Equal(t, 9.0, horizontalIntersectionInfo2.IntersectionPoint.X)
	assert.Equal(t, 0.5, horizontalIntersectionInfo2.IntersectionPoint.Y)

	verticalTestMap := raycastmap.NewSliceMap([][]int{
		{1, 0, 0, 0, 0, 0, 0, 0, 0, 1}, // Wall cells from y=0 to y=9
	}, 0.0, 0.0, 0.0, wallValueToStructure)

	verticalIntersectionInfo1 := RaycastRay(&Vector{0.5, 4.5}, &Vector{0.0, -1.0}, verticalTestMap)
	verticalIntersectionInfo2 := RaycastRay(&Vector{0.5, 4.5}, &Vector{0.0, 1.0}, verticalTestMap)

	assert.Equal(t, 0.5, verticalIntersectionInfo1.IntersectionPoint.X)
	assert.Equal(t, 1.0, verticalIntersectionInfo1.IntersectionPoint.Y)
	assert.Equal(t, 0.5, verticalIntersectionInfo2.IntersectionPoint.X)
	assert.Equal(t, 9.0, verticalIntersectionInfo2.IntersectionPoint.Y)
}

func TestRaycast(t *testing.T) {
	var observerX = 22.0
	var observerY = 12.0

	pixelColumnInfos := Raycast(80, &Vector{observerX, observerY}, math.Pi, raycastmap.TestMap1)
	heightFactor := 40.0
	for _, pixelColumnInfo := range pixelColumnInfos {
		height := heightFactor / pixelColumnInfo.PerpendicularDistance
		fmt.Printf("%s%s\n", line(int(heightFactor-height/2.0), " "), line(int(height), fmt.Sprintf("%d", 8)))
	}
}

func line(count int, character string) string {
	builder := strings.Builder{}
	for i := 0; i < count; i++ {
		builder.WriteString(character)
		builder.WriteString(character)
	}

	return builder.String()
}
