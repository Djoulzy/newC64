package graphic

type RGB struct {
	R byte
	G byte
	B byte
}

type Driver interface {
	Init(int, int)
	DirectDrawPixel(int, int, RGB)
	DrawPixel(int, int, RGB)
	UpdateFrame()
	Run()
	IOEvents() *KEYPressed
	SetKeyboardLine(*KEYPressed)
	CloseAll()
}
