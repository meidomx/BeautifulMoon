package glutils

import "github.com/go-gl/gl/v3.3-core/gl"

type Bindable interface {
	Bind()
	ResetBinding()
}

func NewBindableFunc(bindable Bindable, f func()) func() {
	return func() {
		bindable.Bind()
		f()
		bindable.ResetBinding()
	}
}

func GetOpenGLVersion() string {
	version := gl.GoStr(gl.GetString(gl.VERSION))
	return version
}
