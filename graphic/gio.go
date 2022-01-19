package graphic

import (
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
)

const windowTitle = "Go C64 Emulator"

type GIODriver struct {
	winHeight int
	winWidth  int
	window    *app.Window
	gtx       layout.Context
	screen    []byte
	isReady   bool
}

func (G *GIODriver) Init(winWidth, winHeight int) {
	G.isReady = false
	go func() {
		w := app.NewWindow(
			app.Size(unit.Px(float32(winWidth)), unit.Px(float32(winHeight))),
			app.Title(windowTitle),
		)
		if err := G.mainWindowLoop(w); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
}

func (G *GIODriver) mainWindowLoop(w *app.Window) error {
	G.window = w

	var ops op.Ops
	for {
		event := <-w.Events()
		switch evt := event.(type) {
		case system.DestroyEvent:
			return evt.Err
		case system.FrameEvent:
			G.gtx = layout.NewContext(&ops, evt)
			evt.Frame(G.gtx.Ops)
			G.isReady = true
		}
	}
}