package cia6526

func (C *CIA) Read(addr uint16) byte {
	reg := addr - ((addr >> 4) << 4)
	switch reg {
	case ICR:
		val := C.Reg[ICR]
		C.Reg[ICR] = 0
		*C.Signal_Pin = 0
		return val
	default:
		return C.Reg[reg]
	}
}
