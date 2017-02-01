package main

import (
	"fmt"
	goimage "image"
	"math"
	"runtime"
	"unsafe"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"

	"github.com/meidomx/BeautifulMoon"
	"github.com/meidomx/BeautifulMoon/config"
	"github.com/meidomx/BeautifulMoon/display/glutils"
	"github.com/meidomx/BeautifulMoon/engine"
	"github.com/meidomx/BeautifulMoon/engine/api"
	"github.com/meidomx/BeautifulMoon/resource/image"
	"github.com/meidomx/BeautifulMoon/globalutils"
)

func init() {
	// This is needed to arrange that main() runs on main thread.
	// See documentation for functions that are only allowed to be called from the main thread.
	runtime.LockOSThread()
}

type DemoPhaseHandler struct {
}

func (this DemoPhaseHandler) DoInitPhase(ph api.PhaseController, event *api.LoopTriggeredEvent) {
	//fmt.Println("init phase")

	b := make([]api.PhaseTask, 4)
	for i := 0; i < len(b); i++ {
		b[i] = api.PhaseTask{
			Attachment: make([]byte, 512),
		}
	}
	// finally submit task for parallel processing
	ph.BatchSubmitTask(b)
}
func (this DemoPhaseHandler) DoConcurrentPhase(t api.PhaseTask, processorNo int) {
	//fmt.Println("conccurent phase:", processorNo)
}
func (this DemoPhaseHandler) DoFinalPhase(ph api.PhaseController) {
	//fmt.Println("final phase")
}

func main() {
	c := bmoon.NewGlobalConfig()
	c.DisplayConfig.DisplayResolution.Width = 800
	c.DisplayConfig.DisplayResolution.Height = 600
	c.DisplayConfig.WindowResizable = false
	c.DisplayConfig.OpenGLConfig.MAJOR_VERSION = 3
	c.DisplayConfig.OpenGLConfig.MINOR_VERSION = 3
	c.DisplayConfig.FullScreen = false
	c.InternalConfig = config.NewInternalConfig()

	ph := DemoPhaseHandler{}
	en, err := engine.NewAndStartMainLoop(nil, c.InternalConfig, ph)
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
	//glfw.WindowHint(glfw.RefreshRate, 120)

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

	//glfw.SwapInterval(1)

	window.SetInputMode(glfw.CursorMode, glfw.CursorNormal)
	cursor := glfw.CreateStandardCursor(int(glfw.HandCursor))
	window.SetCursor(cursor)

	w, h := window.GetFramebufferSize()
	gl.Viewport(0, 0, int32(w), int32(h))
	texture := getTexture(rgba)

	shaderProgramObject := glutils.NewShaderProgram()
	globalutils.PanicError(shaderProgramObject.SetVertexShader(_VERTEX_SHADER))
	globalutils.PanicError(shaderProgramObject.SetFragmentShader(_FRAGMENT_SHADER))
	globalutils.PanicError(shaderProgramObject.Link())

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

	vertexColorLocation := gl.GetUniformLocation(shaderProgramObject.GetShaderProgramId(), gl.Str("ourColor"+string([]byte{0})))

	last := glfw.GetTime()
	last = 0

	for !window.ShouldClose() {
		// chapter 1
		gl.ClearColor(0.2, 0.3, 0.3, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT)

		t := glfw.GetTime()

		greenValue := float32(math.Sin(float64(t))/2 + 0.5)

		var _ = texture

		shaderProgramObject.UseProgram()
		gl.Uniform4f(vertexColorLocation, 0.0, greenValue, 0.0, 0.0)
		gl.BindVertexArray(vaoId)
		gl.DrawArrays(gl.TRIANGLES, 0, 3)
		gl.BindVertexArray(0)

		shaderProgramObject.UseProgram()
		gl.Uniform4f(vertexColorLocation, 0.0, greenValue, 0.0, 0.0)
		gl.BindVertexArray(eboVAOId)
		gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, eboId)
		gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, unsafe.Pointer(uintptr(0)))
		gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, 0)
		gl.BindVertexArray(0)

		fps := 1 / (t - last)
		window.SetTitle(fmt.Sprintf("%.1f / %.1f", fps, en.GetFPS()))
		last = t

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
    gl_Position = vec4(position.xyz, 1.0);
}
	` + string([]byte{0})
	_FRAGMENT_SHADER = `
#version 330 core

out vec4 color;
uniform vec4 ourColor;

void main()
{
    color = ourColor;
}
	` + string([]byte{0})
)
