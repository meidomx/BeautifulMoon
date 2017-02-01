package glfwutils

import (
	"github.com/go-gl/glfw/v3.2/glfw"

	"github.com/meidomx/BeautifulMoon/config"
)

func InitWindow(c *config.GlobalConfig) (*glfw.Window, error) {
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
		return nil, err
	}

	return window, nil
}
