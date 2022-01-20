package graphic

import (
	"image"
	"image/color"
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
)

const windowTitle = "Go C64 Emulator"

type GIODriver struct {
	winHeight int
	winWidth  int
	window    *app.Window
	gtx       layout.Context
	screen    []color.NRGBA
	isReady   bool
}

func getColor(col RGB) color.NRGBA {
	nrgba := color.NRGBA{}
	nrgba.R = uint8(col.R)
	nrgba.G = uint8(col.G)
	nrgba.B = uint8(col.B)
	nrgba.A = uint8(255)
	return nrgba
}

func (G *GIODriver) Init(winWidth, winHeight int) {
	G.winWidth = winWidth
	G.winHeight = winHeight
	G.screen = make([]color.NRGBA, winWidth*winHeight)
	// G.screen := [winWidth][winHeight]byte
}

func (G *GIODriver) Start() {
	G.isReady = false
	go func() {
		G.window = app.NewWindow(
			app.Size(unit.Px(float32(G.winWidth)), unit.Px(float32(G.winHeight))),
			app.Title(windowTitle),
		)
		if err := G.mainWindowLoop(); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

func (G *GIODriver) mainWindowLoop() error {
	var ops op.Ops
	for {
		event := <-G.window.Events()
		switch evt := event.(type) {
		case system.DestroyEvent:
			return evt.Err
		case system.FrameEvent:
			G.gtx = layout.NewContext(&ops, evt)
			G.apply()
			evt.Frame(G.gtx.Ops)
			G.isReady = true
		}
	}
}

func (G *GIODriver) DrawPixel(x, y int, color RGB) {
	index := y*G.winWidth + x
	G.screen[index] = getColor(color)
	// if G.isReady {
	// 	paint.FillShape(&G.ops, getColor(color), clip.Rect(image.Rect(x, y, x+1, y+1)).Op())
	// }
}

func (G *GIODriver) apply() {
	for y:=0; y < G.winHeight; y++ {
		for x := 0; x < G.winWidth; x++ {
			paint.FillShape(G.gtx.Ops, G.screen[y*G.winHeight+x], clip.Rect(image.Rect(x, y, x+1, y+1)).Op())
		}
	}
}

func (G *GIODriver) UpdateFrame() {
	G.window.Invalidate()
}

func (G *GIODriver) IOEvents() uint {
	defer func() {
		buffer = 0
	}()
	return buffer
}

func (G *GIODriver) CloseAll() {
	// Do nothing
}
