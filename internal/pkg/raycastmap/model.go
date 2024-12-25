package raycastmap

import (
	"github.com/anthonynsimon/bild/blend"
	"image"
)

type Cell struct {
	X, Y      int
	Structure *Structure
}

type Map interface {
	StartX() float64
	StartY() float64
	StartDir() float64
	Width() int
	Height() int
	WallAt(x, y int) bool
	ObstacleAt(x, y int) bool
	StructureAt(x, y int) *Structure
	SpecialAt(x, y int) *Structure
}

type Structure struct {
	Texture  *Texture
	Texture2 *Texture
	Overlay  *Texture

	Item       bool // Something you can pick up (keys, ammo clip, treasures, extra life...)
	Decoration bool // Something that decorates the cell (skeleton bones, large flower pot, bowl of food...)
	Obstacle   bool // Something you cannot move through (wall, large flower pot, floor light... Not doors though)
	Wall       bool // Some cell 100% covered (walls, like doors, walls). Used to render walls.
}

func NewStructure(texture image.Image) *Structure {
	s := &Structure{}
	s.WithTexture(texture)
	s.WithObstacle(true)
	return s
}

func (s *Structure) WithSecondTexture(texture image.Image) *Structure {
	s.Texture2 = NewTexture(texture)
	return s
}

func (s *Structure) WithTexture(texture image.Image) *Structure {
	s.Texture = NewTexture(texture)
	return s
}

func (s *Structure) WithOverlayTexture(overlay image.Image) *Structure {
	s.Overlay = NewTexture(overlay)

	if overlay != nil {
		if s.Texture != nil {
			s.Texture.img = blend.Normal(s.Texture.img, overlay)
		}
		if s.Texture2 != nil {
			s.Texture2.img = blend.Normal(s.Texture2.img, overlay)
		}
	}

	return s
}

func (s *Structure) IsObstacle() bool {
	return s.Obstacle
}

func (s *Structure) WithObstacle(obstacle bool) *Structure {
	s.Obstacle = obstacle
	return s
}

func (s *Structure) IsWall() bool {
	return s.Wall
}

func (s *Structure) WithWall(wall bool) *Structure {
	s.Wall = wall
	return s
}

func (s *Structure) IsItem() bool {
	return s.Item
}

func (s *Structure) WithItem(item bool) *Structure {
	s.Item = item
	return s
}

func (s *Structure) IsDecoration() bool {
	return s.Decoration
}

func (s *Structure) WithDecoration(decoration bool) *Structure {
	s.Decoration = decoration
	return s
}
