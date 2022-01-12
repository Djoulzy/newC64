package vic6569

func (V *VIC) Read(addr uint16) byte {
	reg := addr - ((addr >> 6) << 6)
	return V.Reg[reg]
}
