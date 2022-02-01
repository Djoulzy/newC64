package graphic

import (
	"log"
	"os"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	Xadjust = -75
	Yadjust = -20
)

type KEYPressed struct {
	KeyCode uint
	Mode    uint
}

type SDLDriver struct {
	winHeight int
	winWidth  int
	window    *sdl.Window
	renderer  *sdl.Renderer
	texture   *sdl.Texture
	screen    []byte
	keybLine  *KEYPressed
}

func (S *SDLDriver) DirectDrawPixel(x, y int, color RGB) {
	S.renderer.SetDrawColor(byte(color.R), byte(color.G), byte(color.B), 255)
	S.renderer.DrawPoint(int32(x), int32(y))
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

	sdl.SetHint(sdl.HINT_RENDER_SCALE_QUALITY, "0")

	S.window, S.renderer, err = sdl.CreateWindowAndRenderer(int32(S.winWidth*2), int32(S.winHeight*2), sdl.WINDOW_SHOWN|sdl.WINDOW_RESIZABLE)
	S.window.SetTitle("Go Commodore 64")

	S.texture, err = S.renderer.CreateTexture(sdl.PIXELFORMAT_RGB24, sdl.TEXTUREACCESS_STREAMING, int32(S.winWidth), int32(S.winHeight))
	if err != nil {
		panic(err)
	}

	S.screen = make([]byte, S.winWidth*S.winHeight*3)
}

func (S *SDLDriver) SetKeyboardLine(line *KEYPressed) {
	S.keybLine = line
}

func (S *SDLDriver) UpdateFrame() {

	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch t := event.(type) {
		case *sdl.QuitEvent:
			os.Exit(1)
		case *sdl.KeyboardEvent:
			switch t.Type {
			case sdl.KEYDOWN:
				S.keybLine.KeyCode = uint(t.Keysym.Sym)
				S.keybLine.Mode = 0
				switch t.Keysym.Mod {
				case 1:
					if S.keybLine.KeyCode != sdl.K_LSHIFT {
						S.keybLine.Mode = sdl.K_LSHIFT
					}
				case 2:
					if S.keybLine.KeyCode != sdl.K_RSHIFT {
						S.keybLine.Mode = sdl.K_RSHIFT
					}
				case 3:
					fallthrough
				default:
					S.keybLine.Mode = S.keybLine.KeyCode
				}
				log.Printf("KEY DOWN : %d %d", S.keybLine.KeyCode, S.keybLine.Mode)
			case sdl.KEYUP:
				// *S.keybLine = 1073742049
				S.keybLine.KeyCode = 0
				S.keybLine.Mode = 0
			}
		default:
			// buffer = 0
		}
	}
	// S.renderer.SetDrawColor(0, 0, 0, 255)
	// S.renderer.Clear()
	S.texture.Update(nil, S.screen, S.winWidth*3)
	S.renderer.Copy(S.texture, nil, nil)
	S.renderer.Present()
}

func (S *SDLDriver) IOEvents() *KEYPressed {
	return S.keybLine
}

func (S *SDLDriver) Run() {
	var isOpen bool = true

	for isOpen {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch t := event.(type) {
			case *sdl.QuitEvent:
				isOpen = false
				break
			case *sdl.KeyboardEvent:
				switch t.Type {
				case sdl.KEYDOWN:
					S.keybLine.KeyCode = uint(t.Keysym.Sym)
					S.keybLine.Mode = 0
					switch t.Keysym.Mod {
					case 1:
						if S.keybLine.KeyCode != sdl.K_LSHIFT {
							S.keybLine.Mode = sdl.K_LSHIFT
						}
					case 2:
						if S.keybLine.KeyCode != sdl.K_RSHIFT {
							S.keybLine.Mode = sdl.K_RSHIFT
						}
					case 3:
						S.keybLine.Mode = 0
					}
					log.Printf("KEY DOWN : %d %d", S.keybLine.KeyCode, S.keybLine.Mode)
				case sdl.KEYUP:
					// *S.keybLine = 1073742049
					S.keybLine.KeyCode = 0
					S.keybLine.Mode = 0
				}
			default:
				// buffer = 0
			}
		}

		// S.renderer.SetDrawColor(0, 0, 0, 255)
		// S.renderer.Clear()
		// S.renderer.SetDrawColor(255, 255, 255, 255)
		// for i := 0; i < 1000; i++ {
		// 	S.renderer.DrawPoint(int32(rand.Intn(S.winWidth)), int32(rand.Intn(S.winHeight)))
		// }
		S.renderer.Present()
	}

	S.CloseAll()
}
