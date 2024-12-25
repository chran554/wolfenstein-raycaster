package raycastmap

import (
	"bytes"
	"image/png"
	"maze/internal/pkg/wolf3d"
	"maze/internal/pkg/wolf3d/resources"
	"os"
)

var (
	overlayPath = "internal/pkg/raycastmap/overlay/"

	StructureNone                    = &Structure{}                                            // Nothing, void, "waste of empty space"
	StructureDoor                    = structure2T("WAL00098", "WAL00099").WithObstacle(false) //
	StructureGreyStoneWall1          = structure2T("WAL00000", "WAL00001")                     // Grey stone wall
	StructureGreyStoneWall2          = structure2T("WAL00002", "WAL00003")                     // Grey stone wall
	StructureSwastikaFlagOnStoneWall = structure2TO("WAL00004", "WAL00005")                    // Red Swastika flag on stone wall
	StructureFramedMoronOnStoneWall  = structure2TO("WAL00006", "WAL00007")                    // Hitler frame on grey stone wall
	StructureEmptyCellClosed         = structure2T("WAL00008", "WAL00009")                     // Cell door (closed door, empty cell)
	StructureEagleStoneArch          = structure2TO("WAL00010", "WAL00011")                    // Eagle in grey stone arch
	StructureSkeletonCellClosed      = structure2T("WAL00012", "WAL00013")                     // Cell door (closed door, skeleton in cell)
	StructureBlueStoneWall1          = structure2T("WAL00014", "WAL00015")                     // Blue stone wall
	StructureBlueStoneWall2          = structure2T("WAL00016", "WAL00017")                     // Blue stone wall
	StructureFramedEagleOnWoodWall   = structure2TO("WAL00018", "WAL00019")                    // Eagle frame on wood wall
	StructureFramedMoronOnWodWall    = structure2TO("WAL00020", "WAL00021")                    // Hitler frame on wood wall
	StructureWoodWall                = structure2T("WAL00022", "WAL00023")                     // Wood wall
	StructureExitDoor                = structure2T("WAL00040", "WAL00043")                     // Exit door
	StructureElevatorDoor            = structure2T("WAL00102", "WAL00103").WithObstacle(false) // Elevator-ish(?) door
	StructureUnknown                 = &Structure{Texture: NewTextureFromFile("overlay/question-mark.png")}

	SpecialNone                          = &Structure{}
	SpecialStartPointFacingNorth         = &Structure{}
	SpecialStartPointFacingEast          = &Structure{}
	SpecialStartPointFacingSouth         = &Structure{}
	SpecialStartPointFacingWest          = &Structure{}
	SpecialBluePuddle                    = structureT("SPR00002")                    // Blue puddle on the floor
	SpecialGreenBarrel                   = structureT("SPR00003").WithObstacle(true) // Green barrel
	SpecialWoodTable                     = structureT("SPR00004").WithObstacle(true) // Wood table
	SpecialGreenLampOnFloor              = structureT("SPR00005").WithObstacle(true) // Green lamp on the floor
	SpecialYellowCrystalChandelierInRoof = structureT("SPR00006")                    // Yellow crystal chandelier in the roof
	SpecialWhiteBowlWithFood             = structureT("SPR00008")                    // White bowl with brown food
	SpecialPlantInGoldFlowerPot          = structureT("SPR00010").WithObstacle(true) // Plant in gold pot
	SpecialSkeletonOnFloor               = structureT("SPR00011")                    // Skeleton lying on the floor
	SpecialPlantInBlueFlowerPot          = structureT("SPR00013").WithObstacle(true) // Brown plant in blue pot
	SpecialBlueFlowerPot                 = structureT("SPR00014").WithObstacle(true) // Blue flower pot
	SpecialRoundTable                    = structureT("SPR00015").WithObstacle(true) // Round table
	SpecialGreenLampInRoof               = structureT("SPR00016")                    // Green lamp in the roof
	SpecialKnightArmour                  = structureT("SPR00018").WithObstacle(true) // Knight armour statue
	SpecialHeapOfBones                   = structureT("SPR00021")                    // Heap of bones
	SpecialBrownBowl                     = structureT("SPR00025")                    // Brown bowl
	SpecialChickenDrumSticks             = structureT("SPR00026")                    // Chicken drumstick on plate
	SpecialMedKit                        = structureT("SPR00027")                    // Med-kit
	SpecialAmmoClip                      = structureT("SPR00028")                    // Ammo clip
	SpecialAutomaticRifle                = structureT("SPR00029")                    // Automatic rifle
	SpecialTreasureGoldCross             = structureT("SPR00031")                    // Treasure gold cross
	SpecialTreasureGoldCup               = structureT("SPR00032")                    // Treasure gold cup
	SpecialTreasureChest                 = structureT("SPR00033")                    // Treasure chest
	SpecialBlueOrb                       = structureT("SPR00035")                    // Blue face orb (extra life)
	SpecialBrownBarrel                   = structureT("SPR00037")                    // Brown barrel
	SpecialStoneWellBlueContent          = structureT("SPR00038").WithObstacle(true) // Stone well, blue liquid
	SpecialStoneWellNoContent            = structureT("SPR00039").WithObstacle(true) // Stone well, no liquid (empty)
	SpecialFlagOnPole                    = structureT("SPR00041").WithObstacle(true) // Flag on standing pole
	SpecialBrownGuard                    = structureT("SPR00050")                    // Brown guard
	SpecialBrownDog                      = structureT("SPR00107")                    // Brown dog
	SpecialDeadGuard                     = structureT("SPR00095")                    // Brown guard dead
	SpecialHiddenDoor                    = &Structure{Texture: NewTextureFromFile("overlay/cross.png")}
	SpecialUnknown                       = &Structure{Texture: NewTextureFromFile("overlay/question-mark.png")}
)

type WolfensteinMap struct {
	levelMaps []wolf3d.LevelMap
	level     int
}

func NewWolfensteinMap(level int) (*WolfensteinMap, error) {
	levelMaps, err := wolf3d.Wolfenstein3DMap()
	if err != nil {
		return nil, err
	}
	return &WolfensteinMap{levelMaps: levelMaps, level: level}, nil
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

func (w *WolfensteinMap) ObstacleAt(x, y int) bool {
	return w.StructureAt(x, y).IsObstacle() || w.SpecialAt(x, y).IsObstacle()
}

func (w *WolfensteinMap) WallAt(x, y int) bool {
	return w.StructureAt(x, y).IsWall()
}

func (w *WolfensteinMap) SpecialAt(x, y int) *Structure {
	specialPlane := 1
	specialValue := w.levelMaps[w.level].Value(specialPlane, x, w.Height()-1-y)

	var special *Structure

	switch specialValue {
	case 0x00:
		special = SpecialNone
	case 0x13:
		special = SpecialStartPointFacingNorth
	case 0x14:
		special = SpecialStartPointFacingEast
	case 0x15:
		special = SpecialStartPointFacingSouth
	case 0x16:
		special = SpecialStartPointFacingWest
	case 0x17:
		special = SpecialBluePuddle
	case 0x18:
		special = SpecialGreenBarrel
	case 0x19:
		special = SpecialWoodTable
	case 0x1A:
		special = SpecialGreenLampOnFloor
	case 0x1B:
		special = SpecialYellowCrystalChandelierInRoof
	case 0x1D:
		special = SpecialWhiteBowlWithFood
	case 0x1F:
		special = SpecialPlantInGoldFlowerPot
	case 0x20:
		special = SpecialSkeletonOnFloor
	case 0x22:
		special = SpecialPlantInBlueFlowerPot
	case 0x23:
		special = SpecialBlueFlowerPot
	case 0x24:
		special = SpecialRoundTable
	case 0x25:
		special = SpecialGreenLampInRoof
	case 0x27:
		special = SpecialKnightArmour
	case 0x2A:
		special = SpecialHeapOfBones
	case 0x2E:
		special = SpecialBrownBowl
	case 0x2F:
		special = SpecialChickenDrumSticks
	case 0x30:
		special = SpecialMedKit
	case 0x31:
		special = SpecialAmmoClip
	case 0x32:
		special = SpecialAutomaticRifle
	case 0x34:
		special = SpecialTreasureGoldCross
	case 0x35:
		special = SpecialTreasureGoldCup
	case 0x36:
		special = SpecialTreasureChest
	case 0x38:
		special = SpecialBlueOrb
	case 0x3A:
		special = SpecialBrownBarrel
	case 0x3B:
		special = SpecialStoneWellBlueContent
	case 0x3C:
		special = SpecialStoneWellNoContent
	case 0x3E:
		special = SpecialFlagOnPole
	case 0x5A, 0x5B:
		special = SpecialBrownGuard
	case 0x70, 0x71, 0x72, 0x73:
		special = SpecialBrownGuard
	case 0x7C:
		special = SpecialDeadGuard
	case 0xB4, 0xB5, 0xB6, 0xB7:
		special = SpecialBrownGuard
	case 0xB8, 0xB9, 0xBA, 0xBB:
		special = SpecialBrownGuard // Patrolling guard (hard)
	case 0x90, 0x91, 0x92, 0x93:
		special = SpecialBrownGuard
	case 0x94, 0x95, 0x96, 0x97:
		special = SpecialBrownGuard // Patrolling guard (medium)
	case 0x6C, 0x6D, 0x6E, 0x6F:
		special = SpecialBrownGuard
	case 0x62:
		special = SpecialHiddenDoor
	case 0x86, 0x87, 0x88, 0x89:
		special = SpecialBrownDog // Standing dog
	case 0x8A, 0x8B, 0x8C, 0x8D:
		special = SpecialBrownDog // Patrolling dog
	case 0xAA, 0xAB, 0xAC, 0xAD:
		special = SpecialBrownDog // Standing dog (medium)
	case 0xAE, 0xAF, 0xB0, 0xB1:
		special = SpecialBrownDog // Patrolling dog (medium)
	case 0xCE, 0xCF, 0xD0, 0xD1:
		special = SpecialBrownDog // Standing dog (hard)
	case 0xD2, 0xD3, 0xD4, 0xD5:
		special = SpecialBrownDog // Patrolling dog (hard)
	default:
		special = SpecialUnknown
	}

	return special
}

func (w *WolfensteinMap) StructureAt(x, y int) *Structure {
	wallPlane := 0
	structureValue := w.levelMaps[w.level].Value(wallPlane, x, w.Height()-1-y)

	structure := StructureNone

	switch structureValue {
	case 0x00:
		structure = StructureNone
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
	case 0x6A, 0x6B, 0x6C, 0x6D, 0x6E, 0x6F,
		0x70, 0x71, 0x72, 0x73, 0x74, 0x75, 0x76, 0x77, 0x78, 0x79, 0x7A, 0x7B, 0x7C, 0x7D, 0x7E, 0x7F,
		0x80, 0x81, 0x82, 0x83, 0x84, 0x85, 0x86, 0x87, 0x88, 0x89, 0x8A, 0x8B, 0x8C, 0x8D, 0x8E, 0x8F:
		structure = StructureNone // Floor and/or roof tile?
	default:
		structure = StructureUnknown
	}

	return structure
}

func structureT(t1 string) *Structure {
	return NewStructure(resources.ImageOrPanic(t1)).
		WithObstacle(true).
		WithWall(true)
}

func structure2T(t1, t2 string) *Structure {
	return NewStructure(resources.ImageOrPanic(t1)).
		WithSecondTexture(resources.ImageOrPanic(t2)).
		WithObstacle(true).
		WithWall(true)
}

func structure2TO(t1, t2 string) *Structure {
	overlayData, _ := os.ReadFile(overlayPath + "never-again.png")
	overlayImage, _ := png.Decode(bytes.NewBuffer(overlayData))

	return NewStructure(resources.ImageOrPanic(t1)).
		WithSecondTexture(resources.ImageOrPanic(t2)).
		WithOverlayTexture(overlayImage).
		WithObstacle(true).
		WithWall(true)
}
