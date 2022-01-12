package vic6569

func (V *VIC) SpreadWrite(reg uint16, val byte) {
	var i uint16
	V.Reg[reg] = val
	for i = 0; i < 10; i++ {
		V.io.Val[reg+i*0x40] = val
	}
}

func (V *VIC) Write(reg uint16, val byte) {
	switch reg {
	case REG_CTRL1:
		V.RasterIRQ &= 0x7FFF
		V.RasterIRQ |= uint16(val&RST8) << 8
		V.SpreadWrite(REG_CTRL1, val)
	case REG_RASTER:
		V.RasterIRQ = V.RasterIRQ&0x8000 + uint16(val)
	case REG_IRQ:
		V.SpreadWrite(REG_IRQ, val)
	case REG_SETIRQ:
		V.SpreadWrite(REG_SETIRQ, val)
	}
}

func (V *VIC) testBit(reg uint16, mask byte) bool {
	if V.Reg[reg]&mask == mask {
		return true
	}
	return false
}

// func (V *VIC) registersManagement() {

// 	if V.io.LastAccess[REG_CTRL1] == memory.WRITE || V.io.LastAccess[REG_RASTER] == memory.WRITE {
// 		V.RasterIRQ = uint16(V.io.Val[REG_CTRL1]&0b10000000) << 8
// 		V.RasterIRQ += uint16(V.io.Val[REG_RASTER])
// 		V.io.LastAccess[REG_CTRL1] = memory.NONE
// 		V.io.LastAccess[REG_RASTER] = memory.NONE
// 	}

// 	if V.io.LastAccess[REG_IRQ] == memory.WRITE {
// 		V.io.VicRegWrite(REG_IRQ, V.io.Val[REG_IRQ]&0b01111111, memory.NONE)
// 		// *V.IRQ_Pin = 0
// 	}
// }
