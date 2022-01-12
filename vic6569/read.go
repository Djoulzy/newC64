package vic6569

func (V *VIC) Read(reg uint16) byte {
	return V.Reg[reg]
}
