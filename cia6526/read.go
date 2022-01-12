package cia6526

	// if C.mem[ICR].IsRead {
	// 	C.mem[ICR].Zone[mem.IO] = 0
	// 	*C.Signal_Pin = 0
	// 	C.mem[ICR].IsRead = false
	// }

func (C *CIA) Read(addr uint16) byte {
	reg := addr - ((addr >> 4) << 4)
	return C.Reg[reg]
}