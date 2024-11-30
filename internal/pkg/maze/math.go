package maze

import (
	"fmt"
	"math"
)

type Vector struct {
	X, Y float64
}

func NewVector(x, y float64) *Vector {
	return &Vector{X: x, Y: y}
}

func NewDirectionVector(angle float64) *Vector {
	return &Vector{X: math.Cos(angle), Y: math.Sin(angle)}
}

func (v *Vector) String() string {
	return fmt.Sprintf("{x: %.2f, y: %.2f}", v.X, v.Y)
}

func (v *Vector) Add(vector *Vector) *Vector {
	return &Vector{X: v.X + vector.X, Y: v.Y + vector.Y}
}

func (v *Vector) AddXY(X, Y float64) *Vector {
	return &Vector{X: v.X + X, Y: v.Y + Y}
}

func (v *Vector) Sub(vector *Vector) *Vector {
	return &Vector{X: v.X - vector.X, Y: v.Y - vector.Y}
}

func (v *Vector) Scale(t float64) *Vector {
	return &Vector{X: v.X * t, Y: v.Y * t}
}

func (v *Vector) Normalized() *Vector {
	return v.Scale(1.0 / v.Length())
}

func (v *Vector) Length() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

func (v *Vector) Flip() *Vector {
	return v.Scale(-1.0)
}

func (v *Vector) PerpendicularClockwise() *Vector {
	return &Vector{X: v.Y, Y: -v.X}
}

func (v *Vector) PerpendicularCounterClockwise() *Vector {
	return &Vector{X: -v.Y, Y: v.X}
}

type Line struct {
	Start *Vector
	End   *Vector

	length            float64
	heading           *Vector
	headingNormalized *Vector
}

func NewLine(start, end *Vector) *Line {
	return &Line{Start: start, End: end}
}

func (l *Line) Length() float64 {
	if l.length == 0 {
		l.length = l.Heading().Length()
	}
	return l.length
}

func (l *Line) Heading() *Vector {
	if l.heading == nil {
		heading := l.End.Sub(l.Start)
		l.heading = heading
	}
	return l.heading
}

func (l *Line) NormalizedHeading() *Vector {
	if l.headingNormalized == nil {
		headingNormalized := l.Heading().Normalized()
		l.headingNormalized = headingNormalized
	}
	return l.headingNormalized
}

func (l *Line) SegmentPoint(t float64) *Vector {
	return l.Start.Add(l.Heading().Scale(t))
}

func (l *Line) Point(distance float64) *Vector {
	return l.Start.Add(l.NormalizedHeading().Scale(distance))
}

func (l *Line) SegmentIntersectionPoint(line *Line) (intersection bool, p *Vector) {
	intersection, t := l.LineSegmentIntersection(line)

	if intersection {
		return true, l.SegmentPoint(t)
	}

	return false, nil
}

// LineSegmentIntersection gives the intersection between this line and another line.
// This reports intersection if intersection Vector is within this line segment
// (not on the line segment extension before start Vector and after end Vector)
func (l *Line) LineSegmentIntersection(line *Line) (intersection bool, t float64) {
	intersection, t, _ = LineIntersection(l, line)
	return intersection && t >= 0 && t <= 1, t
}

// LineIntersection gives the intersection between this line and another line.
// This reports intersection if the intersection Vector is anywhere along the extension of the line segments.
//
// The intersection Vector falls within the line1 segment if 0 ≤ t ≤ 1.
// The intersection Vector falls within the line2 segment if 0 ≤ u ≤ 1.
func LineIntersection(line1 *Line, line2 *Line) (intersection bool, t, u float64) {
	// https://en.wikipedia.org/wiki/Line%E2%80%93line_intersection

	// L1 is a line going from Vector P1s to Vector P1e
	// L2 is a line going from Vector P2s to Vector P2e

	// L1 = P1s + t*(P1e-P1s)
	// L2 = P2s + u*(P2e-P2s)

	// v = (P1sx-P1ex)(P2sy-P2ey)-(P1sy-P1ey)(P2sx-P2ex)
	// t =  ( ((P1sx-P2sx)(P2sy-P2ey)-(P1sy-P2sy)(P2sx-P2ex)) / v)
	// u = -( ((P1sx-P1ex)(P1sy-P2sy)-(P1sy-P1ey)(P1sx-P2sx)) / v)

	v := (line1.Start.X-line1.End.X)*(line2.Start.Y-line2.End.Y) - (line1.Start.Y-line1.End.Y)*(line2.Start.X-line2.End.X)
	if v == 0 {
		return false, 0.0, 0.0
	}

	tt := (line1.Start.X-line2.Start.X)*(line2.Start.Y-line2.End.Y) - (line1.Start.Y-line2.Start.Y)*(line2.Start.X-line2.End.X)
	uu := (line1.Start.X-line1.End.X)*(line1.Start.Y-line2.Start.Y) - (line1.Start.Y-line1.End.Y)*(line1.Start.X-line2.Start.X)

	return true, tt / v, -uu / v
}
