package main

import (
	"fmt"
	goimage "image"
	"runtime"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.2/glfw"

	"github.com/meidomx/BeautifulMoon"
	"github.com/meidomx/BeautifulMoon/engine"
	"github.com/meidomx/BeautifulMoon/resource/image"
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
	c.DisplayConfig.OpenGLConfig.MAJOR_VERSION = 2
	c.DisplayConfig.OpenGLConfig.MINOR_VERSION = 1
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

	for !window.ShouldClose() {
		gl.BindTexture(gl.TEXTURE_2D, texture)

		gl.Begin(gl.QUADS)
		gl.TexCoord2f(0.0, 1.0)
		gl.Vertex3f(-1.0, -1.0, 0.0)
		gl.TexCoord2f(1.0, 1.0)
		gl.Vertex3f(1.0, -1.0, 0.0)
		gl.TexCoord2f(1.0, 0.0)
		gl.Vertex3f(1.0, 1.0, 0.0)
		gl.TexCoord2f(0.0, 0.0)
		gl.Vertex3f(-1.0, 1.0, 0.0)
		gl.End()

		// Do OpenGL stuff.
		window.SwapBuffers()
		glfw.PollEvents()
	}
}

func getTexture(rgba *goimage.RGBA) uint32 {
	var texture uint32
	gl.Enable(gl.TEXTURE_2D)
	gl.GenTextures(1, &texture)
	fmt.Println(texture)
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
