package pla906114

import (
	"newC64/cia6526"
	"newC64/memory"
	"newC64/vic6569"
)

type MEM struct {
	Size     int
	readOnly bool
	Cells    []byte
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
	LORAM  = 0b00000001
	HIRAM  = 0b00000010
	CHAREN = 0b00000100
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
	setting       *byte
	startLocation [5]int
	Mem           [5]*memory.MEM

	vic  *vic6569.VIC
	cia1 *cia6526.CIA
	cia2 *cia6526.CIA
}
