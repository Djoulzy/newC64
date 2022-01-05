package graphic

import (
	"os"

	"github.com/veandco/go-sdl2/sdl"
)

var buffer uint

type SDLDriver struct {
	winHeight int
	winWidth  int
	window    *sdl.Window
	renderer  *sdl.Renderer
	texture   *sdl.Texture
	screen    []byte
}

func (S *SDLDriver) DrawPixel(x, y int, color RGB) {
	index := (y*S.winWidth + x) * 3
	S.screen[index] = byte(color.R)
	S.screen[index+1] = byte(color.G)
	S.screen[index+2] = byte(color.B)
}

func (S *SDLDriver) CloseAll() {
	S.window.Destroy()
	S.renderer.Destroy()
	S.texture.Destroy()
	sdl.Quit()
}

func (S *SDLDriver) Init(winWidth, winHeight int) {
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
	S.window.SetResizable(true)
	S.window.SetSize(int32(S.winWidth*2), int32(S.winHeight*2))

	S.renderer, err = sdl.CreateRenderer(S.window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		panic(err)
	}

	S.texture, err = S.renderer.CreateTexture(sdl.PIXELFORMAT_RGB24, sdl.TEXTUREACCESS_STATIC, int32(S.winWidth), int32(S.winHeight))
	if err != nil {
		panic(err)
	}

	S.screen = make([]byte, S.winWidth*S.winHeight*3)
}

func (S *SDLDriver) UpdateFrame() {

	S.texture.Update(nil, S.screen, S.winWidth*3)
	S.renderer.Copy(S.texture, nil, nil)
	S.renderer.Present()

	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch t := event.(type) {
		case *sdl.QuitEvent:
			os.Exit(1)
		case *sdl.KeyboardEvent:
			buffer = uint(t.Keysym.Sym)
			// switch KeyCode {
			// case sdl.K_l:
			// 	fmt.Printf("A")
			// }
		default:
			buffer = 0
		}
	}
}

func (S *SDLDriver) IOEvents() uint {
	defer func() {
		buffer = 0
	}()
	return buffer
	// event := sdl.PollEvent()
	// switch t := event.(type) {
	// case *sdl.QuitEvent:
	// 	os.Exit(1)
	// case *sdl.KeyboardEvent:
	// 	KeyCode := t.Keysym.Sym
	// 	fmt.Printf("%v\n", KeyCode)
	// 	fmt.Printf("%v\n", sdl.K_a)
	// 	switch KeyCode {
	// 	case sdl.K_a:
	// 		fmt.Printf("A")
	// 	}
	// default:
	// }
}
