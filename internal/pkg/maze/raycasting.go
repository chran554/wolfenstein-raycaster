package maze

import (
	"math"
	"maze/internal/pkg/raycastmap"
)

const humongousLarge = 1e30

type IntersectionInfo struct {
	Hit                        bool
	PerpendicularDistance      float64
	Wall                       *raycastmap.Cell
	Side                       int // North-South wall (=0) or East-West wall (=1)
	WallSideIntersectionOffset float64
	ObserverPoint              *Vector
	IntersectionPoint          *Vector
	IntersectionCosAngle       float64
}

func RaycastRay(start *Vector, rayDir *Vector, worldMap raycastmap.Map) (intersectionInfo IntersectionInfo) {
	// Direction Vector is always of length 1.0.
	// The direction Vector points in the direction the observer is viewing along (at the center of observer view).
	//var rayDirX = math.Cos(directionAngle)
	//var rayDirY = math.Sin(directionAngle)

	// dirLengthSqr := rayDirX*rayDirX + rayDirY*rayDirY
	// if dirLengthSqr != 1.0 {
	// 	invDirLength := 1.0 / math.Sqrt(dirLengthSqr)
	// 	rayDirX = rayDirX * invDirLength
	// 	rayDirY = rayDirY * invDirLength
	// }

	// which box of the map we're in
	var mapX = int(start.X)
	var mapY = int(start.Y)

	// length of ray from one x or y-side to next x or y-side
	deltaDistX := humongousLarge
	deltaDistY := humongousLarge
	if rayDir.X != 0.0 {
		deltaDistX = math.Abs(1.0 / rayDir.X) // An optimization that does not make it possible for easy Euclidean distance calculation
	}
	if rayDir.Y != 0.0 {
		deltaDistY = math.Abs(1.0 / rayDir.Y) // An optimization that does not make it possible for easy Euclidean distance calculation
	}

	// what direction to step in x or y-direction (either +1 or -1)
	var stepDirectionX int
	var stepDirectionY int

	// length of ray from current position to next x or y-side
	var sideDistX float64
	var sideDistY float64

	// calculate step and initial sideDist
	if rayDir.X < 0.0 {
		stepDirectionX = -1
		sideDistX = (start.X - float64(mapX)) * deltaDistX
	} else {
		stepDirectionX = 1
		sideDistX = (float64(mapX) + 1.0 - start.X) * deltaDistX
	}

	if rayDir.Y < 0.0 {
		stepDirectionY = -1
		sideDistY = (start.Y - float64(mapY)) * deltaDistY
	} else {
		stepDirectionY = 1
		sideDistY = (float64(mapY) + 1.0 - start.Y) * deltaDistY
	}

	// perform DDA
	side := -1 // was a NS or an EW wall hit?
	hit := false
	for !hit {
		// Jump to the next map square, either in x-direction or in y-direction
		if sideDistX < sideDistY {
			sideDistX += deltaDistX
			mapX += stepDirectionX
			side = 0 // EW side
		} else {
			sideDistY += deltaDistY
			mapY += stepDirectionY
			side = 1 // NS side
		}

		hit = worldMap.WallAt(mapX, mapY) // Check if ray has hit a wall
	}

	// Calculate distance projected on a camera direction
	// (Euclidean distance would give fisheye effect)
	var perpWallDist float64
	if side == 0 {
		perpWallDist = sideDistX - deltaDistX
	} else {
		perpWallDist = sideDistY - deltaDistY
	}

	// wallIntersectionOffset is the intersection vector offset on the wall, range [0.0, 1.0] (from wall start from left to right).
	// Used for texture pixel column selection.
	var wallIntersectionOffset float64
	if side == 0 { // We are hit west or east side of a wall cell.
		if rayDir.X <= 0 {
			// We are facing/tracing westwards, thus we hit the east side of wall cell
			wallIntersectionOffset = start.Y + perpWallDist*rayDir.Y
		} else if rayDir.X > 0 {
			// We are facing/tracing eastwards, thus we hit the west side of wall cell
			wallIntersectionOffset = start.Y + perpWallDist*rayDir.Y
			wallIntersectionOffset = 1.0 - wallIntersectionOffset
		}
	} else if side == 1 { // We are hit nest or south side of a wall cell.
		if rayDir.Y >= 0.0 {
			// We are facing/tracing northwards, thus we hit the south side of wall cell
			wallIntersectionOffset = start.X + perpWallDist*rayDir.X
		} else if rayDir.Y < 0 {
			wallIntersectionOffset = start.X + perpWallDist*rayDir.X
			wallIntersectionOffset = 1.0 - wallIntersectionOffset
		}
	}
	wallIntersectionOffset -= math.Floor(wallIntersectionOffset)

	// Calculate intersection point coordinate
	intersectionPoint := &Vector{X: float64(mapX), Y: float64(mapY)}
	intersectionCosAngle := 1.0
	if side == 0 { // We are hit west or east side of a wall cell.
		if rayDir.X <= 0 {
			// We are facing/tracing westwards, thus we hit the east side of the wall cell
			// Adjust final cell distance as we hit the east side, thus wall coordinate-cell width (1.0)
			intersectionPoint = intersectionPoint.AddXY(1.0, wallIntersectionOffset)
			intersectionCosAngle = math.Abs(rayDir.Normalized().X)
		} else if rayDir.X > 0 {
			// We are facing/tracing eastwards, thus we hit the west side of the wall cell
			// Adjust final cell distance as we hit the east side, thus wall coordinate-cell width (1.0)
			intersectionPoint = intersectionPoint.AddXY(0.0, 1.0-wallIntersectionOffset)
			intersectionCosAngle = math.Abs(intersectionPoint.Sub(start).Normalized().X)
		}
	} else if side == 1 { // We are hit nest or south side of a wall cell.
		if rayDir.Y >= 0.0 {
			// We are facing/tracing northwards, thus we hit the south side of the wall cell
			// Adjust final cell distance as we hit the east side, thus wall coordinate-cell width (1.0)
			intersectionPoint = intersectionPoint.AddXY(wallIntersectionOffset, 0.0)
			intersectionCosAngle = math.Abs(rayDir.Normalized().Y)
		} else if rayDir.Y < 0 {
			// We are facing/tracing southwards, thus we hit the north side of the wall cell
			// Adjust final cell distance as we hit the east side, thus wall coordinate-cell width (1.0)
			intersectionPoint = intersectionPoint.AddXY(1.0-wallIntersectionOffset, 1.0)
			intersectionCosAngle = math.Abs(rayDir.Normalized().Y)
		}
	}

	intersectionInfo = IntersectionInfo{
		Hit:                        hit,
		PerpendicularDistance:      perpWallDist,
		ObserverPoint:              start,
		IntersectionPoint:          intersectionPoint,
		IntersectionCosAngle:       intersectionCosAngle,
		Wall:                       &raycastmap.Cell{X: mapX, Y: mapY, Structure: worldMap.StructureAt(mapX, mapY)},
		Side:                       side,
		WallSideIntersectionOffset: wallIntersectionOffset,
	}

	return intersectionInfo
}

func Raycast(pixelColumnCount int, observer *Vector, viewDirectionAngle float64, worldMap raycastmap.Map) (pixelColumnInfos []IntersectionInfo) {
	// Direction Vector is always of length 1.0.
	// The direction Vector points in the direction the observer is viewing along (at the center of observer view).
	dirX := math.Cos(viewDirectionAngle)
	dirY := math.Sin(viewDirectionAngle)

	// initial camera plane Vector
	// camera plane Vector is always perpendicular to direction Vector
	fov := 0.66
	planeX := fov * dirY
	planeY := fov * -dirX

	for pixelColumn := 0; pixelColumn < pixelColumnCount; pixelColumn++ {
		// calculate ray position and direction
		cameraX := 2.0*float64(pixelColumn)/float64(pixelColumnCount) - 1.0 // camera plane pos [-1, 1]
		rayDir := &Vector{
			X: dirX + planeX*cameraX,
			Y: dirY + planeY*cameraX,
		}

		pixelColumnInfo := RaycastRay(observer, rayDir, worldMap)
		pixelColumnInfos = append(pixelColumnInfos, pixelColumnInfo)
	}

	return pixelColumnInfos
}
