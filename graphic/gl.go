package graphic

import (
	"os"

	"github.com/go-gl/glfw/v3.3/glfw"
)

type GLDriver struct {
	winHeight int
	winWidth  int
	window    *glfw.Window
	// renderer  *glfw.Renderer
	// texture   *glfw.Texture
	screen []byte
}

func (S *GLDriver) SetPixel(index int, c RGB) {

}

func (S *GLDriver) Draw8pixels(x, y int, fg_color, bg_color RGB, value byte) {

}

func (S *GLDriver) CloseAll() {

}

func (S *GLDriver) Init(winWidth, winHeight int) {
	S.winHeight = winHeight
	S.winWidth = winWidth

	err := glfw.Init()
	if err != nil {
		panic(err)
	}

	S.window, err = glfw.CreateWindow(640, 480, "Testing", nil, nil)
	if err != nil {
		panic(err)
	}

	S.window.MakeContextCurrent()
}

func (S *GLDriver) DisplayFrame() {
	if !S.window.ShouldClose() {
		// Do OpenGL stuff.
		S.window.SwapBuffers()
		glfw.PollEvents()
	} else {
		os.Exit(1)
	}
}
