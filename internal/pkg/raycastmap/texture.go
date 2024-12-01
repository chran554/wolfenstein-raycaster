package raycastmap

import (
	"github.com/anthonynsimon/bild/blend"
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
)

type Texture struct {
	img           image.Image
	dominantColor color.Color
}

func NewTextureWithOverlay(imageFilename string, overlayFilename string) *Texture {
	texture := &Texture{}

	textureImage, _ := readImage(imageFilename)
	overlayImage, _ := readImage(overlayFilename)

	texture.img = textureImage

	if overlayImage != nil {
		texture.img = blend.Normal(textureImage, overlayImage)
	}

	if texture.img != nil {
		texture.dominantColor = averageColor(textureImage)
	}

	return texture
}

func NewTexture(imageFilename string) *Texture {
	texture := &Texture{}

	texture.img, _ = readImage(imageFilename)
	if texture.img != nil {
		texture.dominantColor = averageColor(texture.img)
	}

	return texture
}

func averageColor(textureImage image.Image) color.RGBA {
	var r, g, b, a int
	for y := 0; y < textureImage.Bounds().Dy(); y++ {
		for x := 0; x < textureImage.Bounds().Dx(); x++ {
			pr, pg, pb, pa := textureImage.At(x, y).RGBA()
			r += int(pr)
			g += int(pg)
			b += int(pb)
			a += int(pa)
		}
	}

	amountPixels := textureImage.Bounds().Dx() * textureImage.Bounds().Dy()
	dominantColor := color.RGBA{
		R: uint8((r / amountPixels) >> 8),
		G: uint8((g / amountPixels) >> 8),
		B: uint8((b / amountPixels) >> 8),
		A: uint8((a / amountPixels) >> 8),
	}
	return dominantColor
}

func (t *Texture) DominantColor() color.Color {
	return t.dominantColor
}

// ReadScaledPixelColumn is a low level very specific function to read a pixel column from an image.
//
//	The pixel column to be read from the image is located at offset xOffset [0.0, 1.0] from the left in the image.
//	The pixel column read from the image is scaled (by "closest neighbor") to fit in the supplied data slice.
//	The pixel column data from the image is read from top to bottom and stored in the data slice with start at index 0.
//	The pixel color information is stored in the slice as bytes in the order R, G, B, A (RGBA color model) for each pixel.
//
// Thus, the data slice is four times larger than the count pixel to store.
//
//	The image is considered to be in color model RGBA.
func (t *Texture) ReadScaledPixelColumn(xOffset float64, yOffset float64, yLength float64, data []byte) {
	srcImageWidth := t.img.Bounds().Dx()
	srcImageHeight := t.img.Bounds().Dy()
	srcPixelX := int(float64(srcImageWidth) * xOffset)
	destPixelCount := len(data) / 4 // four bytes per pixel R,G,B,A

	for destPixelIndex := 0; destPixelIndex < destPixelCount; destPixelIndex++ {
		progress := float64(destPixelIndex) / float64(destPixelCount)
		srcPixelY := int(float64(srcImageHeight) * (yLength*progress + yOffset))
		if srcPixelY > srcImageHeight {
			srcPixelY = srcImageHeight - 1
		}

		r, g, b, a := t.img.At(srcPixelX, srcPixelY).RGBA()

		data[destPixelIndex*4+0] = byte(r >> 8)
		data[destPixelIndex*4+1] = byte(g >> 8)
		data[destPixelIndex*4+2] = byte(b >> 8)
		data[destPixelIndex*4+3] = byte(a >> 8)
	}
}

func readImage(filename string) (image.Image, error) {
	file, err := os.Open(filename)
	if err != nil {
		log.Printf("could not read image file: %s\n", err.Error())
		return nil, err
	}
	defer file.Close()

	img, err := png.Decode(file)
	if err != nil {
		log.Printf("could not decode image file: %s\n", err.Error())
		return nil, err
	}

	return img, nil
}
