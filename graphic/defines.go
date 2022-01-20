package graphic

type RGB struct {
	R byte
	G byte
	B byte
}

type Driver interface {
	Init(int, int)
	Start()
	DrawPixel(int, int, RGB)
	UpdateFrame()
	IOEvents() uint
	CloseAll()
}

var buffer uint
