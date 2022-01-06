package mos6510

//
const (
	C_mask byte = 0b11111110
	Z_mask byte = 0b11111101
	I_mask byte = 0b11111011
	D_mask byte = 0b11110111
	B_mask byte = 0b11101111

	V_mask byte = 0b10111111
	N_mask byte = 0b01111111
)

type addressing int

const (
	implied addressing = iota
	immediate
	relative
	zeropage
	zeropageX
	zeropageY
	absolute
	absoluteX
	absoluteY
	indirect
	indirectX
	indirectY
)

type instruction struct {
	name   string
	addr   addressing
	bytes  int
	cycles int
	action func()
}

type cpuState int

const (
	idle cpuState = iota
	readInstruction
	readOperLO
	readOperHI
	compute
)

type memory interface {
	Init()
	Clear()
	Read(uint16) byte
	Write(uint16, byte)
}

// CPU :
type CPU struct {
	PC uint16
	SP byte
	A  byte
	X  byte
	Y  byte
	S  byte

	ram        memory
	instStart  uint16
	inst       instruction
	oper       uint16
	cycleCount int
	state      cpuState
}

// Mnemonic :
var mnemonic map[byte]instruction
