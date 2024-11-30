package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"golang.org/x/image/colornames"
	"image"
	"image/color"
	"math"
	"maze/internal/pkg/maze"
	"maze/internal/pkg/opensimplex"
	"maze/internal/pkg/raycastmap"
	"os"
	"strconv"
	"time"
)

const detectWalls = true

var (
	wolfensteinOriginalWidth  = 320
	wolfensteinOriginalHeight = 200
	scaleFactor               = 2

	observerRadius = 0.2

	useTextures      = false
	useAmbientLight  = 0 // Value: 0 == "ambient light off", 1 == "full ambient light", 2 == "dark ambient light"
	useObserverLight = 0 // Value: 0 == "observer light off", 1 == "observer light", 2 == "observer light animation"
)

func main() {
	application := app.New()
	window := application.NewWindow("Maze")

	renderWidth := wolfensteinOriginalWidth * scaleFactor
	renderHeight := wolfensteinOriginalHeight * scaleFactor

	windowWidth := wolfensteinOriginalWidth * scaleFactor
	windowHeight := wolfensteinOriginalHeight * scaleFactor

	// worldMap, err := raycastmap.WolfMap()
	worldMap, err := raycastmap.NewWolfensteinMap(0)
	if err != nil {
		panic(1)
	}

	observer := &maze.Vector{X: worldMap.StartX(), Y: worldMap.StartY()}
	viewDirectionAngle := worldMap.StartDir()

	noiseGenerator := opensimplex.New(100)

	movementLength := 0.2
	turnSpeed := (math.Pi * 2.0) / (3.0 * 30.0) // one 360 turn in 2 seconds (if frame rate is 30)

	keyUpPressed := false
	keyDownPressed := false
	keyLeftPressed := false
	keyRightPressed := false
	keyAltLeftPressed := false
	keyAltRightPressed := false

	if dc, ok := window.Canvas().(desktop.Canvas); ok {
		dc.SetOnKeyDown(func(event *fyne.KeyEvent) {
			if event.Name == fyne.KeyT {
				useTextures = !useTextures
			} else if event.Name == fyne.KeyA {
				useAmbientLight++
				useAmbientLight = useAmbientLight % 3
			} else if event.Name == fyne.KeyO {
				useObserverLight++
				useObserverLight = useObserverLight % 3
			} else if event.Name == fyne.KeyEscape { // Quick quit
				os.Exit(0)
			} else if event.Name == fyne.KeyUp {
				keyUpPressed = true
			} else if event.Name == fyne.KeyDown {
				keyDownPressed = true
			} else if event.Name == fyne.KeyLeft {
				keyLeftPressed = true
			} else if event.Name == fyne.KeyRight {
				keyRightPressed = true
			} else if event.Physical.ScanCode == 58 { // Left Alt
				keyAltLeftPressed = true
			} else if event.Physical.ScanCode == 59 { // Right Alt
				keyAltRightPressed = true
			} else {
				fmt.Println("Unmapped key down: " + string(event.Name) + ", " + strconv.Itoa(event.Physical.ScanCode))
			}
		})

		dc.SetOnKeyUp(func(event *fyne.KeyEvent) {
			if event.Name == fyne.KeyT {
			} else if event.Name == fyne.KeyA {
			} else if event.Name == fyne.KeyO {
			} else if event.Name == fyne.KeyEscape {
			} else if event.Name == fyne.KeyUp {
				keyUpPressed = false
			} else if event.Name == fyne.KeyDown {
				keyDownPressed = false
			} else if event.Name == fyne.KeyLeft {
				keyLeftPressed = false
			} else if event.Name == fyne.KeyRight {
				keyRightPressed = false
			} else if event.Physical.ScanCode == 58 { // Left Alt
				keyAltLeftPressed = false
			} else if event.Physical.ScanCode == 61 { // Right Alt
				keyAltRightPressed = false
			} else {
				fmt.Println("Unmapped key up:   " + string(event.Name) + ", " + strconv.Itoa(event.Physical.ScanCode))
			}
		})
	}

	imageSize := image.Rect(0, 0, renderWidth, renderHeight)
	img := image.NewRGBA(imageSize)

	rayImage := image.NewRGBA(image.Rect(0, 0, 100, 100))
	rayMapCanvas := canvas.NewImageFromImage(rayImage)
	rayMapCanvas.SetMinSize(fyne.NewSize(float32(rayImage.Bounds().Dx()), float32(rayImage.Bounds().Dy())))

	fpsLabel := widget.NewLabel("")
	posLabel := widget.NewLabel("")
	featureLabel := widget.NewLabel("")

	imgCanvas := canvas.NewImageFromImage(img)
	imgCanvas.FillMode = canvas.ImageFillStretch
	imgCanvas.ScaleMode = canvas.ImageScaleFastest
	infoContainer := container.NewVBox(container.NewHBox(rayMapCanvas, container.NewVBox(container.NewHBox(fpsLabel, posLabel, layout.NewSpacer()), container.NewVBox(featureLabel, layout.NewSpacer())), layout.NewSpacer()))
	window.SetContent(container.NewStack(imgCanvas, infoContainer))

	window.Resize(fyne.NewSize(float32(windowWidth), float32(windowHeight)))

	go func() {
		timestamp := time.Now()
		for {
			if keyLeftPressed {
				if keyAltLeftPressed || keyAltRightPressed {
					// Strafe left
					headingDirection := maze.NewDirectionVector(viewDirectionAngle + math.Pi/2.0)
					obstacleInTheWay := isObstacleInTheWay(headingDirection, observer, worldMap, observerRadius, movementLength)
					if !detectWalls || !obstacleInTheWay {
						observer = observer.Add(headingDirection.Scale(movementLength))
					}
				} else {
					// Turn left
					viewDirectionAngle += turnSpeed

					if viewDirectionAngle > math.Pi*2.0 {
						viewDirectionAngle -= math.Pi * 2.0
					}
				}
			}
			if keyRightPressed {
				if keyAltLeftPressed || keyAltRightPressed {
					// Strafe right
					headingDirection := maze.NewDirectionVector(viewDirectionAngle - math.Pi/2.0)
					obstacleInTheWay := isObstacleInTheWay(headingDirection, observer, worldMap, observerRadius, movementLength)
					if !detectWalls || !obstacleInTheWay {
						observer = observer.Add(headingDirection.Scale(movementLength))
					}
				} else {
					// Turn right
					viewDirectionAngle -= turnSpeed
					if viewDirectionAngle < 0.0 {
						viewDirectionAngle += math.Pi * 2.0
					}
				}
			}
			if keyUpPressed {
				headingDirection := maze.NewDirectionVector(viewDirectionAngle)
				obstacleInTheWay := isObstacleInTheWay(headingDirection, observer, worldMap, observerRadius, movementLength)
				if !detectWalls || !obstacleInTheWay {
					observer = observer.Add(headingDirection.Scale(movementLength))
				}
			}
			if keyDownPressed {
				headingDirection := maze.NewDirectionVector(viewDirectionAngle).Flip()

				if detectWalls {
					if !isObstacleInTheWay(headingDirection, observer, worldMap, observerRadius, movementLength) {
						observer = observer.Add(headingDirection.Scale(movementLength))
					}
				} else {
					observer = observer.Add(headingDirection.Scale(movementLength))
				}
			}

			torchFade := 0.0
			if useObserverLight == 2 {
				noiseSpeed := 500.0 // The higher value, the slower fluctuations in noise function
				noisePosition := float64(uint32(time.Now().UnixMilli())) / noiseSpeed
				smoothRandomValues := noiseGenerator.Eval64(noisePosition) // Value range [-1, 1]

				torchFade = (1.0 - smoothRandomValues) / 2.0 // Compress value range [-1, 1] --> [0, 1]
			}

			torchLight := maze.NewColor(1.0, 1.0, 0.9).FadeTo(maze.NewColor(0.6, 0.5, 0.3), torchFade)
			ambientLight := maze.NewColor(0.2, 0.2, 0.3)

			pixelColumnInfos := maze.Raycast(img.Bounds().Dx(), observer, viewDirectionAngle, worldMap)
			clearImage(img, ambientLight, torchLight)
			if useTextures {
				paintImageTexturized(img, pixelColumnInfos, ambientLight, torchLight)
			} else {
				paintImageColorized(img, pixelColumnInfos, ambientLight, torchLight)
			}

			paintRayMap(rayImage, observer, worldMap)
			rayMapCanvas.Refresh()

			now := time.Now()
			duration := now.Sub(timestamp)
			timestamp = now
			fps := 1000.0 / float64(duration.Milliseconds())

			textureString := "OFF"
			if useTextures {
				textureString = "ON"
			}
			ambientString := "OFF"
			if useAmbientLight == 1 {
				ambientString = "FULL"
			} else if useAmbientLight == 2 {
				ambientString = "LOW"
			}
			observerLightString := "OFF"
			if useObserverLight == 1 {
				observerLightString = "ON"
			} else if useObserverLight == 2 {
				observerLightString = "ANIMATED"
			}

			fpsLabel.SetText(fmt.Sprintf("FPS: %.0f", fps))
			posLabel.SetText(fmt.Sprintf("pos: %+v  dir: %.0f", observer, viewDirectionAngle*(180.0/math.Pi)))
			featureLabel.SetText(fmt.Sprintf("[a] ambient light: %s    [o] observer light: %s    [t] texture: %s", ambientString, observerLightString, textureString))

			imgCanvas.Refresh()
			time.Sleep(20 * time.Millisecond)
		}
	}()

	window.ShowAndRun()
}

func isObstacleInTheWay(headingDirection *maze.Vector, observer *maze.Vector, worldMap *raycastmap.WolfensteinMap, observerRadius float64, movementLength float64) bool {
	info := maze.RaycastRay(observer, headingDirection, worldMap)
	dist := info.IntersectionPoint.Sub(observer).Length()
	return (dist - observerRadius) < movementLength
}

func paintRayMap(mapImage *image.RGBA, observer *maze.Vector, m raycastmap.Map) {
	colorBorder := color.NRGBA{R: 128, G: 16, B: 16, A: 128}

	w := mapImage.Bounds().Dx()
	h := mapImage.Bounds().Dy()
	hw := w / 2
	hh := h / 2

	// Clear with black and fully transparent color
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			mapImage.Set(x, y, color.RGBA{})
		}
	}

	xOffset := int(observer.X) - hw/2
	yOffset := int(observer.Y) - hh/2

	// Draw walls
	for iy := 0; iy < hh; iy++ {
		for ix := 0; ix < hw; ix++ {
			mapx := ix + xOffset
			mapy := iy + yOffset

			c := color.Color(color.RGBA{A: 196})
			if m.WallAt(mapx, mapy) {
				c = m.StructureAt(mapx, mapy).Texture.DominantColor()
			}

			mapImage.Set(2*ix, h-2*iy, c)
			mapImage.Set(2*ix+1, h-2*iy, c)
			mapImage.Set(2*ix, h-(2*iy+1), c)
			mapImage.Set(2*ix+1, h-(2*iy+1), c)
		}
	}

	// Observer (in the middle of map)
	mapImage.Set(hw, hh, colornames.White)
	mapImage.Set(hw+1, hh, colornames.White)
	mapImage.Set(hw, hh+1, colornames.White)
	mapImage.Set(hw+1, hh+1, colornames.White)

	// Draw borders
	for x := 0; x < w; x++ {
		mapImage.Set(x, 0, colorBorder)
		mapImage.Set(x, h-1, colorBorder)
	}
	for y := 0; y < h; y++ {
		mapImage.Set(0, y, colorBorder)
		mapImage.Set(w-1, y, colorBorder)
	}
}

func paintImageTexturized(img *image.RGBA, pixelColumnInfos []maze.IntersectionInfo, ambientLight *maze.Color, torchLight *maze.Color) {
	//w := img.Bounds().Dx()
	h := img.Bounds().Dy()

	for x, pixelColumnInfo := range pixelColumnInfos {
		theoreticalPixelColumnHeight := int(float64(h) / pixelColumnInfo.PerpendicularDistance)
		actualPixelColumnHeight := min(h, theoreticalPixelColumnHeight)

		// Draw scaled texture pixel column
		texture := pixelColumnInfo.Wall.Structure.Texture
		if pixelColumnInfo.Side == 1 && pixelColumnInfo.Wall.Structure.Texture2 != nil && useObserverLight == 0 {
			texture = pixelColumnInfo.Wall.Structure.Texture2 // Use darker texture on East-West facing wall sides of a cell
		}

		xOffset := pixelColumnInfo.WallSideIntersectionOffset
		yLength := float64(actualPixelColumnHeight) / float64(theoreticalPixelColumnHeight)
		yOffset := (1.0 - yLength) / 2.0
		scaledPixelData := make([]byte, actualPixelColumnHeight*4)
		texture.ReadScaledPixelColumn(xOffset, yOffset, yLength, scaledPixelData)

		imageYStart := int(float64(h-actualPixelColumnHeight) / 2.0)
		imgDataOffset := img.PixOffset(x, imageYStart)
		for pixelYIndex := 0; pixelYIndex < theoreticalPixelColumnHeight; pixelYIndex++ {
			imageDataIndex := imgDataOffset + pixelYIndex*img.Stride
			//imageDataIndex := img.PixOffset(x, y1+pixelYIndex)

			if imageDataIndex >= 0 && (imageDataIndex+3) < len(img.Pix) {

				r := scaledPixelData[pixelYIndex*4+0]
				g := scaledPixelData[pixelYIndex*4+1]
				b := scaledPixelData[pixelYIndex*4+2]
				pixelColor := maze.NewColorFromByte(r, g, b)

				distanceAttenuation := 1.0
				cosIntersectionAngle := 1.0
				if useObserverLight == 0 {
					torchLight = maze.NewColor(0.0, 0.0, 0.0)
				} else if useObserverLight >= 1 {
					cosIntersectionAngle = pixelColumnInfo.IntersectionCosAngle

					const attenuationFalloff = 15.0 // Attenuation distance falloff distance setting
					distance := pixelColumnInfo.IntersectionPoint.Sub(pixelColumnInfo.ObserverPoint).Length()
					distanceAttenuation = min(1.0, max(0.0, attenuationFalloff/(distance*distance)))
				}

				if useAmbientLight == 0 {
					ambientLight = maze.NewColor(0.0, 0.0, 0.0)
				} else if useAmbientLight == 1 {
					ambientLight = maze.NewColor(1.0, 1.0, 1.0)
				}

				rb, gb, bb := pixelColor.Mul(ambientLight).Add(pixelColor.Mul(torchLight).Scale(cosIntersectionAngle).Scale(distanceAttenuation)).Bytes()

				img.Pix[imageDataIndex+0] = rb
				img.Pix[imageDataIndex+1] = gb
				img.Pix[imageDataIndex+2] = bb
				img.Pix[imageDataIndex+3] = scaledPixelData[pixelYIndex*4+3] // A (alpha)

				// Middle hair "cross"
				if x == len(pixelColumnInfos)/2 {
					img.Pix[imageDataIndex+0] = 255
				}
			}
		}
	}
}

func paintImageColorized(img *image.RGBA, pixelColumnInfos []maze.IntersectionInfo, ambientLight *maze.Color, torchLight *maze.Color) {
	//w := img.Bounds().Dx()
	h := img.Bounds().Dy()

	for x, pixelColumnInfo := range pixelColumnInfos {
		lineHeight := float64(h) / pixelColumnInfo.PerpendicularDistance
		y1 := (float64(h) - lineHeight) / 2.0
		y2 := y1 + lineHeight

		// Color index from what kind of wall plus brightness index from what side of wall
		c := pixelColumnInfo.Wall.Structure.Texture.DominantColor()
		if pixelColumnInfo.Side == 1 {
			c = pixelColumnInfo.Wall.Structure.Texture2.DominantColor()
		}

		cosIntersectionAngle := pixelColumnInfo.IntersectionCosAngle

		const attenuationFalloff = 15.0 // Attenuation distance falloff distance setting
		distance := pixelColumnInfo.IntersectionPoint.Sub(pixelColumnInfo.ObserverPoint).Length()
		distanceAttenuation := min(1.0, max(0.0, attenuationFalloff/(distance*distance)))

		if useAmbientLight == 0 {
			ambientLight = maze.NewColor(0.0, 0.0, 0.0)
		} else if useAmbientLight == 1 {
			ambientLight = maze.NewColor(1.0, 1.0, 1.0)
		}

		if useObserverLight == 0 {
			torchLight = maze.NewColor(0.0, 0.0, 0.0)
		}

		nc := maze.NewColorFromColor(c)
		c = nc.Mul(ambientLight).Add(nc.Mul(torchLight).Scale(cosIntersectionAngle).Scale(distanceAttenuation)).RGBA()

		drawVerticalLine(img, x, int(y1), int(y2), c)
	}
}

func drawVerticalLine(img *image.RGBA, x int, y1 int, y2 int, c color.Color) {
	for y := 0; y < y2-y1; y++ {
		img.Set(x, y1+y, c)
	}
}

// clearImage paints the "background" of the game. That is, the roof and the floor.
func clearImage(img *image.RGBA, ambientLight *maze.Color, torchLight *maze.Color) {
	colorRoof := maze.NewColor(0.23, 0.23, 0.23)
	colorFloor := maze.NewColor(0.42, 0.42, 0.42)

	imageHeight := float64(img.Bounds().Dy())
	halfImageHeight := imageHeight / 2.0

	for y := 0; y < img.Rect.Dy(); y++ {
		c := colorRoof
		if y >= (img.Rect.Dy() / 2) {
			c = colorFloor
		}

		// Calculation of line height for a wall att perpendicular distance: lineHeight := float64(h) / pixelColumnInfo.PerpendicularDistance
		halfWallHeightAtDistance := math.Abs(halfImageHeight - float64(y))
		distance := 1.0 / (halfWallHeightAtDistance * 2.0 / imageHeight)

		const attenuationFalloff = 15.0 // Attenuation distance falloff distance setting
		distanceAttenuation := min(1.0, max(0.0, attenuationFalloff/(distance*distance)))
		cosAngle := maze.NewVector(distance, 1.0).Normalized().Y // Height above ground should really be 0.5 i.e. half a wall height up from the ground, not 1.0 (but 1.0 yields better result)

		if useAmbientLight == 0 {
			ambientLight = maze.NewColor(0.0, 0.0, 0.0)
		} else if useAmbientLight == 1 {
			ambientLight = maze.NewColor(1.0, 1.0, 1.0)
		}

		if useObserverLight == 0 {
			torchLight = maze.NewColor(0.0, 0.0, 0.0)
		}

		c = c.Mul(ambientLight).Add(c.Mul(torchLight).Scale(cosAngle).Scale(distanceAttenuation))
		rgba := c.RGBA()

		for x := 0; x < img.Rect.Dx(); x++ {
			img.Set(x, y, rgba)
		}
	}
}
