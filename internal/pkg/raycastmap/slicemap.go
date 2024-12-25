package raycastmap

type SliceMap struct {
	mapData            [][]int
	startX, startY     float64
	startDir           float64
	valueToStructureFn func(int) *Structure
}

// NewSliceMap encapsulates map data in the form of a 2D slice of int.
// Data is structured as cell[x][y]. The structured as columns of map cell data.
//
//	data := [][]int{
//	  {1, 2, 3},
//	  {4, 5, 6},
//	  {7, 8, 9},
//	}
//
//	y
//	^
//	| 3 6 9
//	| 2 5 8
//	| 1 4 7
//	+-------> x
func NewSliceMap(mapData [][]int, startX, startY float64, startDir float64, valueToStructureFn func(int) *Structure) SliceMap {
	return SliceMap{mapData: mapData, startX: startX, startY: startY, startDir: startDir, valueToStructureFn: valueToStructureFn}
}

func (sm SliceMap) StructureAt(x, y int) *Structure {
	return sm.valueToStructureFn(sm.mapData[x][y])
}

func (sm SliceMap) SpecialAt(_, _ int) *Structure {
	return &Structure{}
}

func (sm SliceMap) WallAt(x, y int) bool {
	return sm.StructureAt(x, y) != StructureNone
}

func (sm SliceMap) ObstacleAt(x, y int) bool {
	return sm.WallAt(x, y) && sm.StructureAt(x, y) != StructureDoor
}

func (sm SliceMap) Width() int {
	return len(sm.mapData)
}
func (sm SliceMap) Height() int {
	return len(sm.mapData[0])
}
func (sm SliceMap) StartX() float64 {
	return sm.startX
}
func (sm SliceMap) StartY() float64 {
	return sm.startY
}

func (sm SliceMap) StartDir() float64 {
	return sm.startDir
}
