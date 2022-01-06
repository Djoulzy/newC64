package mos6510

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
	COLDSTART_Vector = 0xFFFC
	IRQBRK_Vector    = 0xFFFE
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

type ram interface {
	Init()
	Clear()
	Load(string)
	Read(uint16) byte
	Write(uint16, byte)
	GetView(int, int) interface{}
	Dump(uint16)
}

// CPU :
type CPU struct {
	PC uint16
	SP byte
	A  byte
	X  byte
	Y  byte
	S  byte

	ram        ram
	stack      ram
	instStart  uint16
	instDump   string
	inst       instruction
	oper       uint16
	cycleCount int
	state      cpuState
}

// Mnemonic :
var mnemonic map[byte]instruction
