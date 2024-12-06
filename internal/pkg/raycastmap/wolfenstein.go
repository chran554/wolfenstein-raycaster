package raycastmap

import (
	"maze/internal/pkg/wolf3d"
)

type WolfensteinMap struct {
	levelMaps []wolf3d.LevelMap
	level     int
	plane     int
}

func NewWolfensteinMap(level int) (*WolfensteinMap, error) {
	levelMaps, err := wolf3d.Wolfenstein3DMap()
	if err != nil {
		return nil, err
	}
	return &WolfensteinMap{levelMaps: levelMaps, level: level, plane: 0}, nil
}

func (w *WolfensteinMap) StartX() float64 {
	return 30.5
}

func (w *WolfensteinMap) StartY() float64 {
	return 6.5 // 58.0
}

func (w *WolfensteinMap) StartDir() float64 {
	return 0.0
}

func (w *WolfensteinMap) Width() int {
	return w.levelMaps[w.level].Width
}

func (w *WolfensteinMap) Height() int {
	return w.levelMaps[w.level].Height
}

func (w *WolfensteinMap) WallAt(x, y int) bool {
	return w.StructureAt(x, y).IsWall()
}

func (w *WolfensteinMap) ObstacleAt(x, y int) bool {
	return w.StructureAt(x, y).IsObstacle()
}

func (w *WolfensteinMap) StructureAt(x, y int) Structure {
	wolfStructure := w.levelMaps[w.level].Value(w.plane, x, w.Height()-1-y)

	structure := StructureNone

	switch wolfStructure {
	case 0x01:
		structure = StructureGreyStoneWall1
	case 0x02:
		structure = StructureGreyStoneWall2
	case 0x03:
		structure = StructureSwastikaFlagOnStoneWall
	case 0x04:
		structure = StructureFramedMoronOnStoneWall
	case 0x05:
		structure = StructureEmptyCellClosed
	case 0x06:
		structure = StructureEagleStoneArch
	case 0x07:
		structure = StructureSkeletonCellClosed
	case 0x08:
		structure = StructureBlueStoneWall1
	case 0x09:
		structure = StructureBlueStoneWall2
	case 0x0A:
		structure = StructureFramedEagleOnWoodWall
	case 0x0B:
		structure = StructureFramedMoronOnWodWall
	case 0x0C:
		structure = StructureWoodWall
	case 0x15:
		structure = StructureExitDoor
	case 0x5A:
		structure = StructureDoor
	case 0x5B:
		structure = StructureDoor
	case 0x64:
		structure = StructureElevatorDoor
	}

	return structure
}
