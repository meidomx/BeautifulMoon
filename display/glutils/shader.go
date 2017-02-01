package glutils

import (
	"errors"
	"reflect"
	"unsafe"

	"github.com/go-gl/gl/v3.3-core/gl"
)

type ShaderProgram struct {
	id uint32

	vertexShaderId   uint32
	fragmentShaderId uint32
}

func ResetShaderProgram() {
	gl.UseProgram(0)
}

func NewShaderProgram() *ShaderProgram {
	sp := new(ShaderProgram)
	return sp
}

func (this *ShaderProgram) GetShaderProgramId() uint32 {
	return this.id
}

func (this *ShaderProgram) IsAvailable() bool {
	return this.id > 0
}

func (this ShaderProgram) UseProgram() {
	gl.UseProgram(this.id)
}

func (this *ShaderProgram) Link() error {
	var shaderProgram = gl.CreateProgram()
	gl.AttachShader(shaderProgram, this.vertexShaderId)
	gl.AttachShader(shaderProgram, this.fragmentShaderId)
	gl.LinkProgram(shaderProgram)
	var compileResult int32
	gl.GetProgramiv(shaderProgram, gl.LINK_STATUS, &compileResult)
	if compileResult == 0 {
		data := make([]uint8, 512)
		header := (*reflect.SliceHeader)(unsafe.Pointer(&data))
		gl.GetProgramInfoLog(shaderProgram, 512, nil, (*uint8)(unsafe.Pointer(header.Data)))
		strdata := make([]byte, 512)
		i := 0
		for ; i < 512 && data[i] > 0; i++ {
			strdata[i] = byte(data[i])
		}
		return errors.New(string(strdata[:i]))
	}
	this.id = shaderProgram
	gl.DeleteShader(this.vertexShaderId)
	gl.DeleteShader(this.fragmentShaderId)
	return nil
}

func (this *ShaderProgram) SetVertexShader(src ...string) error {
	var vertexShaderId = gl.CreateShader(gl.VERTEX_SHADER)
	var vertexShaderStr, vertexShaderStrFreeFunc = gl.Strs(src...)
	defer vertexShaderStrFreeFunc()
	gl.ShaderSource(vertexShaderId, int32(len(src)), vertexShaderStr, nil)
	gl.CompileShader(vertexShaderId)
	var compileResult int32
	gl.GetShaderiv(vertexShaderId, gl.COMPILE_STATUS, &compileResult)
	if compileResult == 0 {
		data := make([]uint8, 512)
		header := (*reflect.SliceHeader)(unsafe.Pointer(&data))
		gl.GetShaderInfoLog(vertexShaderId, 512, nil, (*uint8)(unsafe.Pointer(header.Data)))
		strdata := make([]byte, 512)
		i := 0
		for ; i < 512 && data[i] > 0; i++ {
			strdata[i] = byte(data[i])
		}
		return errors.New(string(strdata[:i]))
	}
	this.vertexShaderId = vertexShaderId
	return nil
}

func (this *ShaderProgram) SetFragmentShader(src ...string) error {
	var fragmentShaderId = gl.CreateShader(gl.FRAGMENT_SHADER)
	var fragShaderStr, fragShaderStrFreeFunc = gl.Strs(src...)
	defer fragShaderStrFreeFunc()
	gl.ShaderSource(fragmentShaderId, int32(len(src)), fragShaderStr, nil)
	gl.CompileShader(fragmentShaderId)
	var compileResult int32
	gl.GetShaderiv(fragmentShaderId, gl.COMPILE_STATUS, &compileResult)
	if compileResult == 0 {
		data := make([]uint8, 512)
		header := (*reflect.SliceHeader)(unsafe.Pointer(&data))
		gl.GetShaderInfoLog(fragmentShaderId, 512, nil, (*uint8)(unsafe.Pointer(header.Data)))
		strdata := make([]byte, 512)
		i := 0
		for ; i < 512 && data[i] > 0; i++ {
			strdata[i] = byte(data[i])
		}
		return errors.New(string(strdata[:i]))
	}
	this.fragmentShaderId = fragmentShaderId
	return nil
}
