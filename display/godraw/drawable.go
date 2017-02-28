package godraw

import "image"

type Drawable interface {
	Image() image.Image
	Pos() image.Rectangle
	CanDraw() bool
}
