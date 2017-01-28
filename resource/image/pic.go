package image

import (
	"image"
	"os"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
)

func NewImageFromFile(path string) (image.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	m, _, err := image.Decode(file)
	return m, err
}
