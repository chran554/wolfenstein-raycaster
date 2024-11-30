package maze

import "maze/internal/pkg/raycastmap"

type Creature struct {
	location    *Vector
	viewHeading float64
}

type World struct {
	stage *raycastmap.Map
	hero  *Creature
}
