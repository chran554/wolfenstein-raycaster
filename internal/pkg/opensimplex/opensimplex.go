package opensimplex

import (
	orig "github.com/ojrac/opensimplex-go"
)

// Generator holds parameter values for OpenSimplex noise.
type Generator struct {
	Seed  int64 // Seed holds the seed value for the noise.
	noise orig.Noise
}

// New returns a seeded open simplex noise instance.
func New(seed int64) *Generator {
	return &Generator{
		Seed:  seed,
		noise: orig.New(seed),
	}
}

// Eval64 returns a float64 noise value for the given coordinates.
func (g *Generator) Eval64(dim ...float64) float64 {
	switch len(dim) {
	case 1:
		return g.eval1D64(dim[0])
	case 2:
		return g.eval2D64(dim[0], dim[1])
	case 3:
		return g.eval3D64(dim[0], dim[1], dim[2])
	}

	return 0
}

// eval1D64 generates float64 OpenSimplex noise value from 1-dimensional coordinate.
func (g *Generator) eval1D64(x float64) float64 {
	return g.noise.Eval2(x, x)
}

// eval2D64 generates float64 OpenSimplex noise value from 2-dimensional coordinates.
func (g *Generator) eval2D64(x, y float64) float64 {
	return g.noise.Eval2(x, y)
}

// eval3D64 generates float64 OpenSimplex noise value from 3-dimensional coordinates.
func (g *Generator) eval3D64(x, y, z float64) float64 {
	return g.noise.Eval3(x, y, z)
}
