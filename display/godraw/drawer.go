package godraw

import (
	"image"
	"image/color"
	"image/draw"
	"sync"
)

const bufferSize int = 2

type Drawer struct {
	swapLock sync.RWMutex

	w, h int

	curImageToDraw *image.RGBA
	imgIdx         int
	images         []*image.RGBA
}

func NewDrawer(w, h int) *Drawer {
	d := new(Drawer)
	d.w = w
	d.h = h
	d.imgIdx = 0
	d.images = make([]*image.RGBA, bufferSize)
	for i := 0; i < len(d.images); i++ {
		img := image.NewRGBA(image.Rectangle{Max: image.Point{X: w, Y: h}})
		resetImage(w, h, img)
		d.images[i] = img
	}
	d.curImageToDraw = d.images[d.imgIdx%bufferSize]

	return d
}

func resetImage(w, h int, img *image.RGBA) {
	img.Set(w, h, color.Black)
}

func (this *Drawer) DrawFrame(drawables []Drawable) {
	this.swapLock.RLock()
	dst := this.curImageToDraw
	for i := 0; i < len(drawables); i++ {
		drawable := drawables[i]
		if drawable.CanDraw() {
			srcImg := drawable.Image()
			draw.Draw(dst, drawable.Pos(), srcImg, srcImg.Bounds().Min, draw.Src)
		}
	}
	this.swapLock.RUnlock()
}

func (this *Drawer) Swap() image.Image {
	this.swapLock.Lock()
	finishedImage := this.curImageToDraw
	this.imgIdx++
	this.curImageToDraw = this.images[this.imgIdx%bufferSize]
	resetImage(this.w, this.h, this.curImageToDraw)
	this.swapLock.Unlock()
	return finishedImage
}
