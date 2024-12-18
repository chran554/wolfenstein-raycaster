package raycastmap

type Cell struct {
	X, Y      int
	Structure Structure
}

type Map interface {
	StartX() float64
	StartY() float64
	StartDir() float64
	Width() int
	Height() int
	WallAt(x, y int) bool
	ObstacleAt(x, y int) bool
	StructureAt(x, y int) Structure
}

type Structure struct {
	Texture  *Texture
	Texture2 *Texture
	Overlay  *Texture
}

func (s Structure) IsWall() bool {
	return s != StructureNone
}

func (s Structure) IsObstacle() bool {
	return s.IsWall() && s != StructureDoor
}

/*
00000 Grey stone wall
00002 Grey stone wall
00004 Red Swastika flag on stone wall
00006 Hitler frame on grey stone wall
00007 cell door (closed door, skeleton cell)
00008 cell door (closed door, empty cell)
00010 Eagle in grey stone arch
00014 Blue stone wall
00016 Blue stone wall
00018 Eagle frame on wood wall
00020 Hitler frame on wood wall
00023 Wood wall
00041 Exit door
00102 Elevator-ish(?) door
00098 Cyan door, handle on left
*/

var (
	texturePath = "internal/pkg/wolf3d/resources/extracted/"
	overlayPath = "internal/pkg/raycastmap/overlay/"

	StructureNone                    = Structure{}                                                                                                                                                                                      // Nothing, void, "waste of empty space"
	StructureDoor                    = Structure{Texture: NewTexture(texturePath + "WAL00098.png"), Texture2: NewTexture(texturePath + "WAL00099.png")}                                                                                 // Nothing, void, "waste of empty space"
	StructureGreyStoneWall1          = Structure{Texture: NewTexture(texturePath + "WAL00000.png"), Texture2: NewTexture(texturePath + "WAL00001.png")}                                                                                 // Grey stone wall
	StructureGreyStoneWall2          = Structure{Texture: NewTexture(texturePath + "WAL00002.png"), Texture2: NewTexture(texturePath + "WAL00003.png")}                                                                                 // Grey stone wall
	StructureSwastikaFlagOnStoneWall = Structure{Texture: NewTextureWithOverlay(texturePath+"WAL00004.png", overlayPath+"never-again.png"), Texture2: NewTextureWithOverlay(texturePath+"WAL00005.png", overlayPath+"never-again.png")} // Red Swastika flag on stone wall
	StructureFramedMoronOnStoneWall  = Structure{Texture: NewTextureWithOverlay(texturePath+"WAL00006.png", overlayPath+"never-again.png"), Texture2: NewTextureWithOverlay(texturePath+"WAL00007.png", overlayPath+"never-again.png")} // Hitler frame on grey stone wall
	StructureEmptyCellClosed         = Structure{Texture: NewTexture(texturePath + "WAL00008.png"), Texture2: NewTexture(texturePath + "WAL00009.png")}                                                                                 // Cell door (closed door, empty cell)
	StructureEagleStoneArch          = Structure{Texture: NewTextureWithOverlay(texturePath+"WAL00010.png", overlayPath+"never-again.png"), Texture2: NewTextureWithOverlay(texturePath+"WAL00011.png", overlayPath+"never-again.png")} // Eagle in grey stone arch
	StructureSkeletonCellClosed      = Structure{Texture: NewTexture(texturePath + "WAL00012.png"), Texture2: NewTexture(texturePath + "WAL00013.png")}                                                                                 // Cell door (closed door, skeleton in cell)
	StructureBlueStoneWall1          = Structure{Texture: NewTexture(texturePath + "WAL00014.png"), Texture2: NewTexture(texturePath + "WAL00015.png")}                                                                                 // Blue stone wall
	StructureBlueStoneWall2          = Structure{Texture: NewTexture(texturePath + "WAL00016.png"), Texture2: NewTexture(texturePath + "WAL00017.png")}                                                                                 // Blue stone wall
	StructureFramedEagleOnWoodWall   = Structure{Texture: NewTextureWithOverlay(texturePath+"WAL00018.png", overlayPath+"never-again.png"), Texture2: NewTextureWithOverlay(texturePath+"WAL00019.png", overlayPath+"never-again.png")} // Eagle frame on wood wall
	StructureFramedMoronOnWodWall    = Structure{Texture: NewTextureWithOverlay(texturePath+"WAL00020.png", overlayPath+"never-again.png"), Texture2: NewTextureWithOverlay(texturePath+"WAL00021.png", overlayPath+"never-again.png")} // Hitler frame on wood wall
	StructureWoodWall                = Structure{Texture: NewTexture(texturePath + "WAL00022.png"), Texture2: NewTexture(texturePath + "WAL00023.png")}                                                                                 // Wood wall
	StructureExitDoor                = Structure{Texture: NewTexture(texturePath + "WAL00041.png")}                                                                                                                                     // Exit door
	StructureElevatorDoor            = Structure{Texture: NewTexture(texturePath + "WAL00102.png"), Texture2: NewTexture(texturePath + "WAL00103.png")}                                                                                 // Elevator-ish(?) door
)
