package vic6569

// func (V *VIC) SpreadWrite(reg uint16, val byte) {
// 	var i uint16
// 	V.Reg[reg] = val
// 	for i = 0; i < 10; i++ {
// 		V.io.Val[reg+i*0x40] = val
// 	}
// }

func (V *VIC) Write(addr uint16, val byte) {

	reg := addr - ((addr >> 6) << 6)

	switch reg {
	case REG_X_SPRT_0:
		fallthrough
	case REG_Y_SPRT_0:
		fallthrough
	case REG_X_SPRT_1:
		fallthrough
	case REG_Y_SPRT_1:
		fallthrough
	case REG_X_SPRT_2:
		fallthrough
	case REG_Y_SPRT_2:
		fallthrough
	case REG_X_SPRT_3:
		fallthrough
	case REG_Y_SPRT_3:
		fallthrough
	case REG_X_SPRT_4:
		fallthrough
	case REG_Y_SPRT_4:
		fallthrough
	case REG_X_SPRT_5:
		fallthrough
	case REG_Y_SPRT_5:
		fallthrough
	case REG_X_SPRT_6:
		fallthrough
	case REG_Y_SPRT_6:
		fallthrough
	case REG_X_SPRT_7:
		fallthrough
	case REG_Y_SPRT_7:
		fallthrough
	case REG_MSBS_X_COOR:
		V.Reg[reg] = val
	case REG_CTRL1:
		V.RasterIRQ &= 0x7FFF
		V.RasterIRQ |= uint16(val&RST8) << 8
		V.Reg[REG_CTRL1] = val
	case REG_RASTER:
		V.RasterIRQ = V.RasterIRQ&0x8000 + uint16(val)
	case REG_LP_X:
		fallthrough
	case REG_LP_Y:
		fallthrough
	case REG_SPRT_ENABLED:
		fallthrough
	case REG_CTRL2:
		fallthrough
	case REG_SPRT_Y_EXP:
		fallthrough
	case REG_MEM_LOC:
		fallthrough
	case REG_IRQ:
		fallthrough
	case REG_IRQ_ENABLED:
		fallthrough
	case REG_SPRT_DATA_PRIORITY:
		fallthrough
	case REG_SPRT_MLTCOLOR:
		fallthrough
	case REG_SPRT_X_EXP:
		fallthrough
	case REG_SPRT_SPRT_COLL:
		fallthrough
	case REG_SPRT_DATA_COLL:
		fallthrough
	case REG_BORDER_COL:
		fallthrough
	case REG_BGCOLOR_0:
		fallthrough
	case REG_BGCOLOR_1:
		fallthrough
	case REG_BGCOLOR_2:
		fallthrough
	case REG_BGCOLOR_3:
		fallthrough
	case REG_SPRT_MLTCOLOR_0:
		fallthrough
	case REG_SPRT_MLTCOLOR_1:
		fallthrough
	case REG_COLOR_SPRT_0:
		fallthrough
	case REG_COLOR_SPRT_1:
		fallthrough
	case REG_COLOR_SPRT_2:
		fallthrough
	case REG_COLOR_SPRT_3:
		fallthrough
	case REG_COLOR_SPRT_4:
		fallthrough
	case REG_COLOR_SPRT_5:
		fallthrough
	case REG_COLOR_SPRT_6:
		fallthrough
	case REG_COLOR_SPRT_7:
		V.Reg[reg] = val
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
