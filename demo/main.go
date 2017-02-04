package main

import (
	"fmt"
	goimage "image"
	"math"
	"os"
	"runtime"
	"runtime/pprof"
	"runtime/trace"
	"unsafe"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"

	"github.com/meidomx/BeautifulMoon"
	"github.com/meidomx/BeautifulMoon/config"
	"github.com/meidomx/BeautifulMoon/display"
	"github.com/meidomx/BeautifulMoon/display/glfwutils"
	"github.com/meidomx/BeautifulMoon/display/glutils"
	"github.com/meidomx/BeautifulMoon/engine"
	"github.com/meidomx/BeautifulMoon/engine/api"
	"github.com/meidomx/BeautifulMoon/globalutils"
	"github.com/meidomx/BeautifulMoon/resource/image"
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
	resoConv := display.NewResolutionConverter(c.DisplayConfig.DisplayResolution.Width, c.DisplayConfig.DisplayResolution.Height)

	img, err := image.NewImageFromFile("resource\\1.png")
	globalutils.PanicError(err)
	rgba, err := image.ToRGBA(img)
	globalutils.PanicError(err)

	err = glfw.Init()
	globalutils.PanicError(err)
	defer glfw.Terminate()

	window, err := glfwutils.InitWindow(c)
	globalutils.PanicError(err)

	window.MakeContextCurrent()

	// esc to exit
	window.SetKeyCallback(func(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
		if key == glfw.KeyEscape && action == glfw.Press {
			window.SetShouldClose(true)
		}
	})

	// Initialize Glow
	globalutils.PanicError(gl.Init())
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
	globalutils.PanicFalse(shaderProgramObject.IsAvailable())

	// chapter 2
	trianglePoint := []float32{
		-0.5, -0.5, 0.0,
		0.5, -0.5, 0.0,
		0.0, 0.5, 0.0,
	}
	vaoObject := glutils.NewVertexArrayObject()
	globalutils.PanicFalse(vaoObject.IsAvailable())
	var bufferId uint32
	gl.GenBuffers(1, &bufferId)
	glutils.NewBindableFunc(vaoObject, func() {
		gl.BindBuffer(gl.ARRAY_BUFFER, bufferId)
		gl.BufferData(gl.ARRAY_BUFFER, 4*len(trianglePoint), gl.Ptr(trianglePoint), gl.STATIC_DRAW)
		gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 3*4, unsafe.Pointer(uintptr(0)))
		gl.EnableVertexAttribArray(0)
		gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	})()
	trianglePointDrawFunc := glutils.NewBindableFunc(vaoObject, func() {
		gl.DrawArrays(gl.TRIANGLES, 0, 3)
	})

	point0x, point0y := resoConv.ConvertToOpenGLCoordinator2d32(100, 550)
	point1x, point1y := resoConv.ConvertToOpenGLCoordinator2d32(100, 470)
	point2x, point2y := resoConv.ConvertToOpenGLCoordinator2d32(200, 550)
	point3x, point3y := resoConv.ConvertToOpenGLCoordinator2d32(200, 470)
	rectanglePoint := []float32{
		point0x, point0y, 0.0,
		point1x, point1y, 0.0,
		point2x, point2y, 0.0,
		point3x, point3y, 0.0,
	}
	rectangleIndices := []uint32{
		0, 1, 2,
		1, 2, 3,
	}
	eboVaoObject := glutils.NewVertexArrayObject()
	globalutils.PanicFalse(eboVaoObject.IsAvailable())
	var eboVboId uint32
	gl.GenBuffers(1, &eboVboId)
	var eboId uint32
	gl.GenBuffers(1, &eboId)
	glutils.NewBindableFunc(eboVaoObject, func() {
		gl.BindBuffer(gl.ARRAY_BUFFER, eboVboId)
		gl.BufferData(gl.ARRAY_BUFFER, 4*len(rectanglePoint), gl.Ptr(rectanglePoint), gl.STATIC_DRAW)
		gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, eboId)
		gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, 4*len(rectangleIndices), gl.Ptr(rectangleIndices), gl.STATIC_DRAW)
		gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 3*4, unsafe.Pointer(uintptr(0)))
		gl.EnableVertexAttribArray(0)
		gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, 0)
		gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	})()

	//gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)

	vertexColorLocation := gl.GetUniformLocation(shaderProgramObject.GetShaderProgramId(), gl.Str("ourColor"+string([]byte{0})))

	last := glfw.GetTime()
	last = 0
	frameCnt := 0

	file, err := os.OpenFile("./trace.out", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, os.FileMode(0644))
	globalutils.PanicError(err)
	globalutils.PanicError(trace.Start(file))
	defer trace.Stop()
	pprofFile, err := os.OpenFile("./pprof.out", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, os.FileMode(0644))
	globalutils.PanicError(err)
	defer pprofFile.Close()
	globalutils.PanicError(pprof.StartCPUProfile(pprofFile))
	defer pprof.StopCPUProfile()

	for !window.ShouldClose() {
		// chapter 1
		gl.ClearColor(0.2, 0.3, 0.3, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT)

		t := glfw.GetTime()

		greenValue := float32(math.Sin(float64(t))/2 + 0.5)

		var _ = texture

		shaderProgramObject.UseProgram()
		gl.Uniform4f(vertexColorLocation, 0.0, greenValue, 0.0, 0.0)
		trianglePointDrawFunc()

		shaderProgramObject.UseProgram()
		gl.Uniform4f(vertexColorLocation, 0.0, greenValue, 0.0, 0.0)
		eboVaoObject.Bind()
		gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, eboId)
		gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, unsafe.Pointer(uintptr(0)))
		gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, 0)
		eboVaoObject.ResetBinding()

		duration := t - last
		frameCnt++
		if duration > 0.4 {
			fps := 1 / duration * float64(frameCnt)
			window.SetTitle(fmt.Sprintf("%.1f / %.1f", fps, en.GetFPS()))
			last = t
			frameCnt = 0
		}

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
