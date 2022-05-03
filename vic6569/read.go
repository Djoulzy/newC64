package vic6569

import "github.com/Djoulzy/Tools/clog"

func (V *VIC) Read(addr uint16) byte {
	reg := addr - ((addr >> 6) << 6)
	clog.Trace("VIC", "Read", "addr: %04X - Reg: %02X (%d)", addr, reg, reg)
	switch reg {
	case REG_MEM_LOC:
		// log.Printf("Read base: %04X", V.Reg[reg])
		return V.Reg[reg]
	case REG_RASTER:
		// log.Printf("RasterIRQ: %04X", V.Reg[reg])
		return V.Reg[reg]
	default:
		return V.Reg[reg]
	}
}
