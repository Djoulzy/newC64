package mos6510

import (
	"newC64/confload"
	"newC64/pla906114"
)

//
const (
	C_mask byte = 0b11111110
	Z_mask byte = 0b11111101
	I_mask byte = 0b11111011
	D_mask byte = 0b11110111
	B_mask byte = 0b11101111
	U_mask byte = 0b11011111
	V_mask byte = 0b10111111
	N_mask byte = 0b01111111

	StackStart = 0x0100

	NMI_Vector       = 0xFFFA
	COLDSTART_Vector = 0xFFFC // Go to 0xFCE2
	IRQBRK_Vector    = 0xFFFE
)

var regString [8]string = [8]string{"C", "Z", "I", "D", "B", "U", "V", "N"}

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
	Branching
	CrossPage
)

type instruction struct {
	name   string
	addr   addressing
	bytes  int
	Cycles int
	action func()
}

type cpuState int

const (
	Idle cpuState = iota
	ReadInstruction
	ReadOperLO
	ReadOperHI
	ReadZP
	ReadZP_XY
	ReadAbsolute
	ReadAbsXY
	ReadIndirect
	ReadIndXY_LO
	ReadIndXY_HI
	Compute
	IRQ1
	IRQ2
	IRQ3
	IRQ4
	IRQ5
	IRQ6
	IRQ7
	NMI1
	NMI2
	NMI3
	NMI4
	NMI5
	NMI6
	NMI7
)

// CPU :
type CPU struct {
	PC      uint16
	SP      byte
	A       byte
	X       byte
	Y       byte
	S       byte
	IRQ_pin int
	NMI_pin int

	conf      *confload.ConfigData
	ram       *pla906114.PLA
	stack     []byte
	InstStart uint16
	instDump  string
	instCode  byte
	Inst      instruction

	oper         uint16
	cross_oper   uint16
	val_zp_lo    byte
	val_zp_hi    byte
	val_absolute byte
	val_absXY    byte
	comp_result  byte

	cycleCount  int
	State       cpuState
	ClockCycles *uint16

	NMI_Raised bool
	IRQ_Raised bool
	INT_delay  bool
}

// Mnemonic :
var mnemonic map[byte]instruction
