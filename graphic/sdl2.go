package graphic

import (
	"github.com/veandco/go-sdl2/sdl"
)

type SDL2Driver struct {
	winHeight int
	winWidth  int
	window    *sdl.Window
	renderer  *sdl.Renderer
}

func (S *SDL2Driver) DrawPixel(x, y int, color RGB) {
	S.renderer.SetDrawColor(byte(color.R), byte(color.R), byte(color.R), 255)
	S.renderer.DrawPoint(int32(x), int32(y))
	// S.renderer.Present()
}

func (S *SDL2Driver) CloseAll() {
	S.window.Destroy()
	S.renderer.Destroy()
	sdl.Quit()
}

func (S *SDL2Driver) Init(winWidth, winHeight int) {
	S.winHeight = winHeight
	S.winWidth = winWidth

	err := sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil {
		panic(err)
	}

	sdl.SetHint(sdl.HINT_RENDER_SCALE_QUALITY, "1")

	S.window, err = sdl.CreateWindow("VIC-II", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, int32(S.winWidth), int32(S.winHeight), sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}

	S.renderer, err = sdl.CreateRenderer(S.window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		panic(err)
	}
}

func (S *SDL2Driver) DisplayFrame() {
	S.renderer.Present()
	sdl.PollEvent()
	// for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
	// 	switch event.(type) {
	// 	case *sdl.QuitEvent:
	// 		os.Exit(1)
	// 	}
	// }

}
