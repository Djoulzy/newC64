package pla906114

type memory interface {
	Init()
	Clear()
	Load(string)
	Read(uint16) byte
	Write(uint16, byte)
	GetView(int, int) interface{}
}

type MemType int

const (
	RAM MemType = iota
	KERNAL
	BASIC
	CHAR
	IO
	CART_LO
	CART_HI
)

const (
	StackStart  = 0x0100
	StackEnd    = 0x01FF
	ScreenStart = 0x0400
	ScreenEnd   = 0x07FF
	IOStart     = 0xD000
	CharStart   = 0xD000
	CharEnd     = 0xDFFF
	ColorStart  = 0xD800
	ColorEnd    = 0xDBFF
	IntAddr     = 0xFFFA
	ResetAddr   = 0xFFFC
	BrkAddr     = 0xFFFE
	KernalStart = 0xE000
	KernalEnd   = 0xFFFF
	BasicStart  = 0xA000
	BasicEnd    = 0xC000
	Vic2        = 0x4000
	Vic3        = 0x8000
	Vic4        = 0xC000
)

// RAM :
type PLA struct {
	setting       byte
	startLocation [5]int
	mem           [5]memory
}
