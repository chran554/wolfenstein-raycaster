package raycastmap

import (
	"github.com/stretchr/testify/assert"
	"image"
	"image/color"
	"image/png"
	"os"
	"testing"
)

func TestReadScaledPixelColumn(t *testing.T) {
	t.Skip("Skip scale texture test as it produce a png-file on file system.")

	const bytesPerPixel = 4

	texture := NewTexture("testdata/WAL00000.png")

	dstWidth := 1024
	dstHeight := 1024

	dstImage := image.NewRGBA(image.Rect(0, 0, dstWidth, dstHeight))
	scaledPixelData := make([]byte, dstHeight*bytesPerPixel)

	for dstX := 0; dstX < dstWidth; dstX++ {
		srcImageXOffset := float64(dstX) / float64(dstWidth)
		texture.ReadScaledPixelColumn(srcImageXOffset, 0, 0, scaledPixelData)

		for dstY := 0; dstY < dstHeight; dstY++ {
			r := scaledPixelData[dstY*bytesPerPixel+0]
			g := scaledPixelData[dstY*bytesPerPixel+1]
			b := scaledPixelData[dstY*bytesPerPixel+2]
			a := scaledPixelData[dstY*bytesPerPixel+3]
			dstImage.Set(dstX, dstY, color.RGBA{R: r, G: g, B: b, A: a})
		}
	}

	pngImageFile, err := os.Create("scaled_texture.png")
	assert.NoError(t, err)
	defer pngImageFile.Close()
	err = png.Encode(pngImageFile, dstImage)
	assert.NoError(t, err)
}
