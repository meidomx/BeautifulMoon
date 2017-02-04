package glutils

import (
	"github.com/go-gl/gl/v3.3-core/gl"
)

type VertexArrayObject struct {
	id uint32
}

func ResetVertexArrayObject() {
	gl.BindVertexArray(0)
}

func NewVertexArrayObject() *VertexArrayObject {
	vao := new(VertexArrayObject)
	gl.GenVertexArrays(1, &vao.id)
	return vao
}

func (this VertexArrayObject) GetVaoId() uint32 {
	return this.id
}

func (this VertexArrayObject) Bind() {
	gl.BindVertexArray(this.id)
}

func (this VertexArrayObject) IsAvailable() bool {
	return this.id > 0
}

func (this VertexArrayObject) ResetBinding() {
	ResetVertexArrayObject()
}
