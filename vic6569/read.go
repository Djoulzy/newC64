package vic6569

func (V *VIC) Read(addr uint16) byte {
	reg := addr - ((addr >> 6) << 6)
	switch reg {
	case REG_MEM_LOC:
		// log.Printf("Read base: %04X", V.Reg[reg])
		return V.Reg[reg]
	default:
		return V.Reg[reg]
	}
}
