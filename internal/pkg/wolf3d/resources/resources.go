package resources

import (
	"bytes"
	"embed"
	"fmt"
	"image"
	"image/png"
)

//go:embed extracted
var resources embed.FS

func ImageResource(filename string) (image.Image, error) {
	imageFileData, err := resources.ReadFile(fmt.Sprintf("%s.png", filename))
	if err != nil {
		return nil, err
	}
	return png.Decode(bytes.NewBuffer(imageFileData))
}

func ImageOrPanic(filename string) image.Image {
	imageFileData, err := resources.ReadFile(fmt.Sprintf("extracted/%s.png", filename))
	if err != nil {
		panic(err)
	}

	img, err := png.Decode(bytes.NewBuffer(imageFileData))
	if err != nil {
		panic(err)
	}

	return img
}
