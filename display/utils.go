package display

type ResolutionConverter struct {
	w          int
	h          int
	widthHalf  int
	heightHalf int
}

func NewResolutionConverter(screenX, screenY int) ResolutionConverter {
	return ResolutionConverter{
		w:          screenX,
		h:          screenY,
		widthHalf:  screenX / 2,
		heightHalf: screenY / 2,
	}
}

func (this ResolutionConverter) ConvertToOpenGLCoordinator2d32(x, y int) (float32, float32) {
	return float32(x)/float32(this.widthHalf) - 1.0, 1.0 - float32(y)/float32(this.heightHalf)
}
