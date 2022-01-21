package cia6526

func (C *CIA) Read(addr uint16) byte {
	reg := addr - ((addr >> 4) << 4)
	switch reg {
	case PRA:
		// val := byte(0b00100000)
		C.Reg[PRA] = 0b10000001
		return C.Reg[PRA]
	case PRB:
		// val := byte(0b00001000)
		C.Reg[PRB] = 0b11111111
		return C.Reg[PRB]
	case ICR:
		val := C.Reg[ICR]
		C.Reg[ICR] = 0
		*C.Signal_Pin = 0
		return val
	default:
		return C.Reg[reg]
	}
}
