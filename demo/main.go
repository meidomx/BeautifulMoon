package main

import (
	"fmt"
	goimage "image"
	"runtime"

	"github.com/go-gl/gl/v3.2-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"

	"github.com/meidomx/BeautifulMoon"
	"github.com/meidomx/BeautifulMoon/engine"
	"github.com/meidomx/BeautifulMoon/resource/image"
	"reflect"
	"unsafe"
)

func init() {
	// This is needed to arrange that main() runs on main thread.
	// See documentation for functions that are only allowed to be called from the main thread.
	runtime.LockOSThread()
}

func main() {
	c := bmoon.NewGlobalConfig()
	c.DisplayConfig.DisplayResolution.Width = 800
	c.DisplayConfig.DisplayResolution.Height = 600
	c.DisplayConfig.WindowResizable = false
	c.DisplayConfig.OpenGLConfig.MAJOR_VERSION = 3
	c.DisplayConfig.OpenGLConfig.MINOR_VERSION = 3
	c.DisplayConfig.FullScreen = false

	en, err := engine.NewAndStartMainLoop(nil)
	if err != nil {
		panic(err)
	}
	en.DoNothing()

	//==========================================================
	img, err := image.NewImageFromFile("resource\\1.jpg")
	if err != nil {
		panic(err)
	}
	rgba, err := image.ToRGBA(img)
	if err != nil {
		panic(err)
	}

	err = glfw.Init()
	if err != nil {
		panic(err)
	}
	defer glfw.Terminate()

	if c.DisplayConfig.WindowResizable {
		glfw.WindowHint(glfw.Resizable, glfw.True)
	} else {
		glfw.WindowHint(glfw.Resizable, glfw.False)
	}
	glfw.WindowHint(glfw.ContextVersionMajor, c.DisplayConfig.OpenGLConfig.MAJOR_VERSION)
	glfw.WindowHint(glfw.ContextVersionMinor, c.DisplayConfig.OpenGLConfig.MINOR_VERSION)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	var monitor *glfw.Monitor = nil
	if c.DisplayConfig.FullScreen {
		monitor = glfw.GetPrimaryMonitor()
	}
	window, err := glfw.CreateWindow(c.DisplayConfig.DisplayResolution.Width, c.DisplayConfig.DisplayResolution.Height, "Testing", monitor, nil)
	if err != nil {
		panic(err)
	}

	window.MakeContextCurrent()

	// esc to exit
	window.SetKeyCallback(func(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
		if key == glfw.KeyEscape && action == glfw.Press {
			window.SetShouldClose(true)
		}
	})

	// Initialize Glow
	if err := gl.Init(); err != nil {
		panic(err)
	}
	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Println("OpenGL version", version)

	glfw.SwapInterval(1)

	window.SetInputMode(glfw.CursorMode, glfw.CursorNormal)
	cursor := glfw.CreateStandardCursor(int(glfw.HandCursor))
	window.SetCursor(cursor)

	w, h := window.GetFramebufferSize()
	gl.Viewport(0, 0, int32(w), int32(h))
	texture := getTexture(rgba)

	var vertextShaderId = gl.CreateShader(gl.VERTEX_SHADER)
	var vertexShaderStr, vertexShaderStrFreeFunc = gl.Strs(_VERTEX_SHADER)
	defer vertexShaderStrFreeFunc()
	gl.ShaderSource(vertextShaderId, 1, vertexShaderStr, nil)
	gl.CompileShader(vertextShaderId)
	var compileResult int32
	gl.GetShaderiv(vertextShaderId, gl.COMPILE_STATUS, &compileResult)
	fmt.Println("vertex shader compile:", compileResult, ",", vertextShaderId)
	if compileResult == 0 {
		data := make([]uint8, 512)
		header := (*reflect.SliceHeader)(unsafe.Pointer(&data))
		gl.GetShaderInfoLog(vertextShaderId, 512, nil, (*uint8)(unsafe.Pointer(header.Data)))
		strdata := make([]byte, 512)
		for i := 0; i < 512; i++ {
			strdata[i] = byte(data[i])
		}
		panic(string(strdata))
	}
	var fragmentShaderid = gl.CreateShader(gl.FRAGMENT_SHADER)
	var fragShaderStr, fragShaderStrFreeFunc = gl.Strs(_FRAGMENT_SHADER)
	defer fragShaderStrFreeFunc()
	gl.ShaderSource(fragmentShaderid, 1, fragShaderStr, nil)
	gl.CompileShader(fragmentShaderid)
	gl.GetShaderiv(fragmentShaderid, gl.COMPILE_STATUS, &compileResult)
	fmt.Println("fragment shader compile:", compileResult, ",", fragmentShaderid)
	if compileResult == 0 {
		data := make([]uint8, 512)
		header := (*reflect.SliceHeader)(unsafe.Pointer(&data))
		gl.GetShaderInfoLog(fragmentShaderid, 512, nil, (*uint8)(unsafe.Pointer(header.Data)))
		strdata := make([]byte, 512)
		for i := 0; i < 512; i++ {
			strdata[i] = byte(data[i])
		}
		panic(string(strdata))
	}
	var shaderProgram = gl.CreateProgram()
	gl.AttachShader(shaderProgram, vertextShaderId)
	gl.AttachShader(shaderProgram, fragmentShaderid)
	gl.LinkProgram(shaderProgram)
	gl.GetProgramiv(shaderProgram, gl.LINK_STATUS, &compileResult)
	fmt.Println("program link:", compileResult, ",", shaderProgram)
	if compileResult == 0 {
		data := make([]uint8, 512)
		header := (*reflect.SliceHeader)(unsafe.Pointer(&data))
		gl.GetProgramInfoLog(shaderProgram, 512, nil, (*uint8)(unsafe.Pointer(header.Data)))
		strdata := make([]byte, 512)
		for i := 0; i < 512; i++ {
			strdata[i] = byte(data[i])
		}
		panic(string(strdata))
	}
	gl.DeleteShader(vertextShaderId)
	gl.DeleteShader(fragmentShaderid)

	var vaoId uint32
	gl.GenVertexArrays(1, &vaoId)
	var bufferId uint32
	gl.GenBuffers(1, &bufferId)
	// chapter 2
	trianglePoint := []float32{
		-0.5, -0.5, 0.0,
		0.5, -0.5, 0.0,
		0.0, 0.5, 0.0,
	}
	gl.BindVertexArray(vaoId)
	gl.BindBuffer(gl.ARRAY_BUFFER, bufferId)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(trianglePoint), gl.Ptr(trianglePoint), gl.STATIC_DRAW)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 3*4, unsafe.Pointer(uintptr(0)))
	gl.EnableVertexAttribArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.BindVertexArray(0)

	rectanglePoint := []float32{
		-0.8, -0.8, 0.0,
		-0.8, -0.6, 0.0,
		0.8, -0.8, 0.0,
		0.8, -0.6, 0.0,
	}
	rectangleIndices := []uint32{
		0, 1, 2,
		1, 2, 3,
	}

	var eboVAOId uint32
	gl.GenVertexArrays(1, &eboVAOId)
	var eboVboId uint32
	gl.GenBuffers(1, &eboVboId)
	var eboId uint32
	gl.GenBuffers(1, &eboId)
	gl.BindVertexArray(eboVAOId)
	gl.BindBuffer(gl.ARRAY_BUFFER, eboVboId)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(rectanglePoint), gl.Ptr(rectanglePoint), gl.STATIC_DRAW)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, eboId)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, 4*len(rectangleIndices), gl.Ptr(rectangleIndices), gl.STATIC_DRAW)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 3*4, unsafe.Pointer(uintptr(0)))
	gl.EnableVertexAttribArray(0)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, 0)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.BindVertexArray(0)

	//gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)

	for !window.ShouldClose() {
		// chapter 1
		gl.ClearColor(0.2, 0.3, 0.3, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT)

		var _ = texture

		gl.UseProgram(shaderProgram)
		gl.BindVertexArray(vaoId)
		gl.DrawArrays(gl.TRIANGLES, 0, 3)
		gl.BindVertexArray(0)

		gl.UseProgram(shaderProgram)
		gl.BindVertexArray(eboVAOId)
		gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, eboId)
		gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, unsafe.Pointer(uintptr(0)))
		gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, 0)
		gl.BindVertexArray(0)

		// Do OpenGL stuff.
		window.SwapBuffers()
		glfw.PollEvents()
	}
}

func getTexture(rgba *goimage.RGBA) uint32 {
	var texture uint32
	gl.Enable(gl.TEXTURE_2D)
	gl.GenTextures(1, &texture)
	gl.BindTexture(gl.TEXTURE_2D, texture)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.RGBA,
		int32(rgba.Rect.Size().X),
		int32(rgba.Rect.Size().Y),
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		gl.Ptr(rgba.Pix))
	return texture
}

var (
	_VERTEX_SHADER = `
#version 330 core

layout (location = 0) in vec3 position;

void main()
{
    gl_Position = vec4(position.x, position.y, position.z, 1.0);
}
	` + string([]byte{0})
	_FRAGMENT_SHADER = `
#version 330 core

out vec4 color;

void main()
{
    color = vec4(1.0f, 0.5f, 0.2f, 1.0f);
}
	` + string([]byte{0})
)
