package pla906114

type memory interface {
	Init()
	Clear()
	Load(string)
	Read(uint16) byte
	Write(uint16, byte)
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
	memorySize  = 65536
	stackStart  = 0x0100
	stackEnd    = 0x01FF
	screenStart = 0x0400
	screenEnd   = 0x07FF
	charStart   = 0xD000
	charEnd     = 0xDFFF
	colorStart  = 0xD800
	colorEnd    = 0xDBFF
	intAddr     = 0xFFFA
	resetAddr   = 0xFFFC
	brkAddr     = 0xFFFE
	KernalStart = 0xE000
	KernalEnd   = 0xFFFF
	BasicStart  = 0xA000
	BasicEnd    = 0xC000
	vic2        = 0x4000
	vic3        = 0x8000
	vic4        = 0xC000
)

// RAM :
type PLA struct {
	setting byte
	mem     [4]memory
}
