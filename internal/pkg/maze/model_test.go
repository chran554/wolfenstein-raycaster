package maze

import (
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
)

func TestLine_Length(t *testing.T) {
	line := &Line{
		Start: &Vector{0.0, 0.0},
		End:   &Vector{1.0, 1.0},
	}

	assert.InEpsilon(t, math.Sqrt(2), line.Length(), 0.000001)
}

func TestLine_Intersection(t *testing.T) {
	type TestCase struct {
		name         string
		l1, l2       *Line
		intersection bool
		s            float64
		u            float64
	}

	testCases := []TestCase{
		{name: "intersect in both lines middle", l1: NewLine(&Vector{0.0, 0.0}, &Vector{1.0, 1.0}), l2: NewLine(&Vector{1.0, 0.0}, &Vector{0.0, 1.0}), intersection: true, s: 0.5, u: 0.5},
		{name: "end points intersect", l1: NewLine(&Vector{0.0, 0.0}, &Vector{1.0, 1.0}), l2: NewLine(&Vector{1.0, 0.0}, &Vector{1.0, 1.0}), intersection: true, s: 1.0, u: 1.0},
		{name: "end Vector and start Vector intersect", l1: NewLine(&Vector{0.0, 0.0}, &Vector{1.0, 1.0}), l2: NewLine(&Vector{1.0, 1.0}, &Vector{1.0, 0.0}), intersection: true, s: 1.0, u: 0.0},
		{name: "parallel lines do not intersect", l1: NewLine(&Vector{0.0, 0.0}, &Vector{1.0, 1.0}), l2: NewLine(&Vector{0.0, 1.0}, &Vector{1.0, 2.0}), intersection: false, s: 0.0, u: 0.0},
		{name: "intersect outside segments", l1: NewLine(&Vector{0.0, 0.0}, &Vector{0.0, 1.0}), l2: NewLine(&Vector{0.0, 2.0}, &Vector{1.0, 2.0}), intersection: true, s: 2.0, u: 0.0},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			intersection, s, u := LineIntersection(testCase.l1, testCase.l2)
			assert.Equal(t, testCase.intersection, intersection)
			assert.Equal(t, testCase.s, s)
			assert.Equal(t, testCase.u, u)
		})
	}
}
