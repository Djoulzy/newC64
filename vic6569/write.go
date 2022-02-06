package vic6569

import "log"

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
		V.MODE = (V.MODE & 0b10011111) | val&0b01100000
		V.RasterIRQ &= 0x7FFF
		V.RasterIRQ |= uint16(val&RST8) << 8
		if val&ECM == ECM {
			V.ECM = true
			log.Printf("Graphic mode: %08b", V.MODE)
		}
		if val&BMM == BMM {
			V.BMM = true
			log.Printf("Graphic mode: %08b", V.MODE)
		}
		V.Reg[REG_CTRL1] |= val & 0b01111111
	case REG_RASTER:
		V.RasterIRQ = V.RasterIRQ&0x8000 + uint16(val)
	case REG_LP_X:
		fallthrough
	case REG_LP_Y:
		fallthrough
	case REG_SPRT_ENABLED:
		fallthrough
	case REG_CTRL2:
		V.MODE = (V.MODE & 0b11101111) | val&0b00010000
		if val&MCM == MCM {
			V.MCM = true
			log.Printf("Graphic mode: %08b", V.MODE)
		}
		V.Reg[reg] = val
	case REG_SPRT_Y_EXP:
		V.Reg[reg] = val
	case REG_MEM_LOC:
		V.ScreenBase = uint16(val&0b11110000) << 6
		V.CharBase = uint16(val&0b00001110) << 10
		V.Reg[reg] = val
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
