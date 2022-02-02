package vic6569

import (
	"fmt"
	"newC64/confload"
	"newC64/graphic"
	"newC64/memory"
)

const (
	// Durée d'un cycle PAL : 1.015uS
	screenWidthPAL  = 504
	screenHeightPAL = 312
	rasterWidthPAL  = 403
	rasterHeightPAL = 284
	cyclesPerLine   = 63

	rasterTime = 1                  // Nb of cycle to put 1 byte on a line
	rasterLine = rasterWidthPAL / 8 // Nb of cycle to draw a full line
	fullRaster = rasterLine * rasterHeightPAL

	// lineRefresh   = cyclesPerLine * cpuCycle                   // Time for a line in ms
	// screenRefresh = screenHeightPAL * cyclesPerLine * cpuCycle // Time for a full screen display in ms
	// fps           = 1 / screenRefresh

	winWidth      = screenWidthPAL
	winHeight     = screenHeightPAL
	visibleWidth  = 320
	visibleHeight = 200

	firstVBlankLine  = 300
	lastVBlankLine   = 15
	firstDisplayLine = 51
	lastDisplayLine  = 250

	// firstHBlankCol  = 453
	// lastHBlankCol   = 50
	// visibleFirstCol = 92
	// visibleLastCol  = 412
)

func (V *VIC) Init(ram *memory.MEM, io *memory.MEM, chargen *memory.MEM, video interface{}, conf *confload.ConfigData) {
	V.graph = video.(graphic.Driver)
	V.graph.Init(winWidth, winHeight)
	V.conf = conf

	// V.io = io.GetView(0, 0x0400)
	V.color = io.GetView(colorStart, 1024)
	// V.screen = ram.GetView(screenStart, 1024)

	V.bankMem[3].Init(2, 0x4000)
	V.bankMem[3].Attach("RAM", 0, 0, ram.Val[0x0000:0x4000])
	V.bankMem[3].Attach("Char ROM", 1, 1, chargen.Val)
	V.bankMem[3].Show()

	V.bankMem[2].Init(1, 0x4000)
	V.bankMem[2].Attach("RAM", 0, 0, ram.Val[0x4000:0x8000])
	V.bankMem[2].Show()

	V.bankMem[1].Init(2, 0x4000)
	V.bankMem[1].Attach("RAM", 0, 0, ram.Val[0x8000:0xC000])
	V.bankMem[1].Attach("Char ROM", 1, 1, chargen.Val)
	V.bankMem[1].Show()

	V.bankMem[0].Init(1, 0x4000)
	V.bankMem[0].Attach("RAM", 0, 0, ram.Val[0xC000:])
	V.bankMem[0].Show()

	// V.bankMem[3].Val = append(append(append(V.bankMem[3].Val, ram.Val[0x0000:0x1000]...), chargen.Val...), ram.Val[0x2000:0x4000]...)
	// V.bankMem[2].Val = ram.Val[0x4000:0x8000]
	// V.bankMem[1].Val = append(append(ram.Val[0x8000:0x9000], chargen.Val...), ram.Val[0xA000:0xC000]...)
	// V.bankMem[0].Val = ram.Val[0xC000:]

	V.BA = true
	V.VCBASE = 0
	V.BeamX = 0
	V.BeamY = 0
	V.cycle = 1
	V.RasterIRQ = 0xFFFF
	V.SystemClock = 0
	V.BankSel = 3
}

func (V *VIC) Disassemble() string {
	var buf string
	buf = fmt.Sprintf("RstX: %03d - RstY: %03d - RC: %02d - VC: %03X - VCBase: %03X - VMLI: %02d", V.BeamX, V.BeamY, V.RC, V.VC, V.VCBASE, V.VMLI)
	return buf
}

func (V *VIC) saveRasterPos(val int) {
	V.Reg[REG_RASTER] = byte(val)
	mask := byte(val>>1) & RST8
	res := V.Reg[REG_CTRL1] & 0b01111111
	V.Reg[REG_CTRL1] = res | mask
}

func (V *VIC) readVideoMatrix() {
	if !V.BA {
		V.ColorBuffer[V.VMLI] = V.color.Val[V.VC] & 0b00001111
		V.CharBuffer[V.VMLI] = V.bankMem[V.BankSel].Read(V.ScreenBase+V.VC)
		// fmt.Printf("VMLI: %02X - VC: %02X - Screen Code: %d - Color: %04X\n", V.VMLI, V.VC, V.CharBuffer[V.VMLI], V.ColorBuffer[V.VMLI])
	}
}

func (V *VIC) drawChar(X int, Y int) {
	if V.drawArea && (V.Reg[REG_CTRL1]&DEN > 0) {
		charAddr := (uint16(V.CharBuffer[V.VMLI]) << 3) + uint16(V.RC)
		charData := V.bankMem[V.BankSel].Read(V.CharBase+charAddr)
		// fmt.Printf("SC: %02X - RC: %d - %04X - %02X = %08b\n", V.CharBuffer[V.VMLI], V.RC, charAddr, charData, charData)
		// if V.CharBuffer[V.VMLI] == 0 {
		// 	fmt.Printf("Raster: %d - Cycle: %d - BA: %t - VMLI: %d - VCBASE/VC: %d/%d - RC: %d - Char: %02X\n", Y, X, V.BA, V.VMLI, V.VCBASE, V.VC, V.RC, V.CharBuffer[V.VMLI])
		// }
		for column := 0; column < 8; column++ {
			bit := byte(0b10000000 >> column)
			if charData&bit > 0 {
				V.graph.DrawPixel(X+column, Y, Colors[V.ColorBuffer[V.VMLI]])
			} else {
				V.graph.DrawPixel(X+column, Y, Colors[V.Reg[REG_BGCOLOR_0]&0b00001111])
			}
		}
		V.VMLI++
		V.VC++
	} else if V.visibleArea {
		for column := 0; column < 8; column++ {
			V.graph.DrawPixel(X+column, Y, Colors[V.Reg[REG_BORDER_COL]&0b00001111])
		}
	}
}

func (V *VIC) Run() bool {
	V.SystemClock++
	V.saveRasterPos(V.BeamY)

	V.visibleArea = (V.BeamY > lastVBlankLine) && (V.BeamY < firstVBlankLine)
	// V.displayArea = (V.BeamY >= firstDisplayLine) && (V.BeamY <= lastDisplayLine) && V.visibleArea
	V.displayArea = (V.BeamY >= firstDisplayLine) && (V.BeamY <= lastDisplayLine)
	V.BeamX = (V.cycle - 1) * 8
	V.drawArea = ((V.cycle > 15) && (V.cycle < 56)) && V.displayArea

	V.BA = !(((V.BeamY-firstDisplayLine)%8 == 0) && V.displayArea && (V.cycle > 11) && (V.cycle < 55))

	// if V.drawArea {
	// 	fmt.Printf("Raster: %d - Cycle: %d - BA: %t - VMLI: %d - VCBASE/VC: %d/%d - RC: %d - Char: %02X\n", V.BeamY, V.cycle, V.BA, V.VMLI, V.VCBASE, V.VC, V.RC, V.CharBuffer[V.VMLI])
	// }

	switch V.cycle {
	case 1:
		if V.testBit(REG_IRQ_ENABLED, IRQ_RST) {
			if V.RasterIRQ == uint16(V.BeamY) {
				//fmt.Printf("\nIRQ: %04X - %04X", V.RasterIRQ, uint16(V.BeamY))
				// fmt.Println("Rastrer Interrupt")
				V.Reg[REG_IRQ] = V.Reg[REG_IRQ] | 0b10000001
				*V.IRQ_Pin = 1
			}
		}
	case 2:
	case 3:
	case 4:
	case 5:
	case 6:
	case 7:
	case 8:
	case 9:
	case 10:
	case 11: // Debut de la zone visible
		V.drawChar(V.BeamX, V.BeamY)
	case 12:
		V.drawChar(V.BeamX, V.BeamY)
	case 13:
		V.drawChar(V.BeamX, V.BeamY)
	case 14:
		V.VC = V.VCBASE
		V.VMLI = 0
		if !V.BA {
			V.RC = 0
		}
		V.drawChar(V.BeamX, V.BeamY)
	case 15: // Debut de la lecture de la memoire video en mode BadLine
		V.drawChar(V.BeamX, V.BeamY)
		V.readVideoMatrix()
	case 16: // Debut de la zone d'affichage
		V.drawChar(V.BeamX, V.BeamY)
		V.readVideoMatrix()
	case 17:
		V.drawChar(V.BeamX, V.BeamY)
		V.readVideoMatrix()
	case 18:
		V.drawChar(V.BeamX, V.BeamY)
		V.readVideoMatrix()
	case 19:
		V.drawChar(V.BeamX, V.BeamY)
		V.readVideoMatrix()
	case 20:
		V.drawChar(V.BeamX, V.BeamY)
		V.readVideoMatrix()
	case 21:
		V.drawChar(V.BeamX, V.BeamY)
		V.readVideoMatrix()
	case 22:
		V.drawChar(V.BeamX, V.BeamY)
		V.readVideoMatrix()
	case 23:
		V.drawChar(V.BeamX, V.BeamY)
		V.readVideoMatrix()
	case 24:
		V.drawChar(V.BeamX, V.BeamY)
		V.readVideoMatrix()
	case 25:
		V.drawChar(V.BeamX, V.BeamY)
		V.readVideoMatrix()
	case 26:
		V.drawChar(V.BeamX, V.BeamY)
		V.readVideoMatrix()
	case 27:
		V.drawChar(V.BeamX, V.BeamY)
		V.readVideoMatrix()
	case 28:
		V.drawChar(V.BeamX, V.BeamY)
		V.readVideoMatrix()
	case 29:
		V.drawChar(V.BeamX, V.BeamY)
		V.readVideoMatrix()
	case 30:
		V.drawChar(V.BeamX, V.BeamY)
		V.readVideoMatrix()
	case 31:
		V.drawChar(V.BeamX, V.BeamY)
		V.readVideoMatrix()
	case 32:
		V.drawChar(V.BeamX, V.BeamY)
		V.readVideoMatrix()
	case 33:
		V.drawChar(V.BeamX, V.BeamY)
		V.readVideoMatrix()
	case 34:
		V.drawChar(V.BeamX, V.BeamY)
		V.readVideoMatrix()
	case 35:
		V.drawChar(V.BeamX, V.BeamY)
		V.readVideoMatrix()
	case 36:
		V.drawChar(V.BeamX, V.BeamY)
		V.readVideoMatrix()
	case 37:
		V.drawChar(V.BeamX, V.BeamY)
		V.readVideoMatrix()
	case 38:
		V.drawChar(V.BeamX, V.BeamY)
		V.readVideoMatrix()
	case 39:
		V.drawChar(V.BeamX, V.BeamY)
		V.readVideoMatrix()
	case 40:
		V.drawChar(V.BeamX, V.BeamY)
		V.readVideoMatrix()
	case 41:
		V.drawChar(V.BeamX, V.BeamY)
		V.readVideoMatrix()
	case 42:
		V.drawChar(V.BeamX, V.BeamY)
		V.readVideoMatrix()
	case 43:
		V.drawChar(V.BeamX, V.BeamY)
		V.readVideoMatrix()
	case 44:
		V.drawChar(V.BeamX, V.BeamY)
		V.readVideoMatrix()
	case 45:
		V.drawChar(V.BeamX, V.BeamY)
		V.readVideoMatrix()
	case 46:
		V.drawChar(V.BeamX, V.BeamY)
		V.readVideoMatrix()
	case 47:
		V.drawChar(V.BeamX, V.BeamY)
		V.readVideoMatrix()
	case 48:
		V.drawChar(V.BeamX, V.BeamY)
		V.readVideoMatrix()
	case 49:
		V.drawChar(V.BeamX, V.BeamY)
		V.readVideoMatrix()
	case 50:
		V.drawChar(V.BeamX, V.BeamY)
		V.readVideoMatrix()
	case 51:
		V.drawChar(V.BeamX, V.BeamY)
		V.readVideoMatrix()
	case 52:
		V.drawChar(V.BeamX, V.BeamY)
		V.readVideoMatrix()
	case 53:
		V.drawChar(V.BeamX, V.BeamY)
		V.readVideoMatrix()
	case 54: // Dernier lecture de la matrice video ram
		V.drawChar(V.BeamX, V.BeamY)
		V.readVideoMatrix()
	case 55: // Fin de la zone de display
		V.drawChar(V.BeamX, V.BeamY)
	case 56: // Debut de la zone visible
		V.drawChar(V.BeamX, V.BeamY)
	case 57:
		V.drawChar(V.BeamX, V.BeamY)
	case 58:
		if V.RC == 7 {
			V.VCBASE = V.VC
		}
		if V.displayArea {
			V.RC++
		}
		V.drawChar(V.BeamX, V.BeamY)
	case 59:
		V.drawChar(V.BeamX, V.BeamY)
	case 60:
	case 61:
	case 62:
	case 63:
	}
	// V.BeamX += 8
	V.cycle++
	if V.cycle > cyclesPerLine {
		V.cycle = 1
		V.BeamY++
		if V.BeamY >= screenHeightPAL {
			V.BeamY = 0
			V.VCBASE = 0
			if V.conf.Globals.Display {
				V.graph.UpdateFrame()
			}
		}
		// if V.conf.Globals.Disassamble == true {
		// 	if V.conf.Globals.Display {
		// 		V.graph.UpdateFrame()
		// 	}
		// }
	}
	return V.BA
}
