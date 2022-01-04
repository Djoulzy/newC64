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
	immediate addressing = iota
	zeropage
	zeropageX
	absolute
	absoluteX
	absoluteY
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

// CPU :
type CPU struct {
	PC uint16
	SP byte
	A  byte
	X  byte
	Y  byte
	S  byte

	inst       instruction
	operLO     byte
	operHI     byte
	cycleCount int
	state      cpuState
}

// Mnemonic :
var mnemonic map[byte]instruction
