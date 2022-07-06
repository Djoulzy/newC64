package vic6569

import (
	"fmt"
	"log"
	"newC64/config"

	"github.com/Djoulzy/emutools/render"

	"github.com/Djoulzy/emutools/mem"
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

	BankStart0 = 0xC000
	BankStart1 = 0x8000
	BankStart2 = 0x4000
	BankStart3 = 0x0000

	DisplayOriginX = 80
	DisplayOriginY = 16
)

func (V *VIC) Init(ram []byte, io []byte, chargen []byte, video *render.SDL2Driver, conf *config.ConfigData) {
	V.graph = video
	V.graph.Init(winWidth-DisplayOriginX-32, winHeight-DisplayOriginY-12, "Go Commodore 64", false, conf.Disassamble)
	V.conf = conf

	V.color = io[colorStart : colorStart+1024]

	V.bankMem = mem.InitBanks(4, &V.BankSel)

	V.bankMem.Layouts[3] = mem.InitConfig(0x4000)
	V.bankMem.Layouts[3].Attach("RAM", 0, ram[BankStart3:BankStart3+0x4000], mem.READWRITE, false)
	V.bankMem.Layouts[3].Attach("Char ROM", 0x1000, chargen, mem.READONLY, false)

	V.bankMem.Layouts[2] = mem.InitConfig(0x4000)
	V.bankMem.Layouts[2].Attach("RAM", 0, ram[BankStart2:BankStart2+0x4000], mem.READWRITE, false)

	V.bankMem.Layouts[1] = mem.InitConfig(0x4000)
	V.bankMem.Layouts[1].Attach("RAM", 0, ram[BankStart1:BankStart1+0x4000], mem.READWRITE, false)
	V.bankMem.Layouts[1].Attach("Char ROM", 0x1000, chargen, mem.READONLY, false)

	V.bankMem.Layouts[0] = mem.InitConfig(0x4000)
	V.bankMem.Layouts[0].Attach("RAM", 0, ram[BankStart0:BankStart0+0x4000], mem.READWRITE, false)

	V.BA = true
	V.VCBASE = 0
	V.BeamX = 0
	V.BeamY = 0
	V.cycle = 1
	V.RasterIRQ = 0xFFFF
	V.SystemClock = 0
	V.MODE = 0
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
	V.ColorBuffer[V.VMLI] = V.color[V.VC] & 0b00001111
	V.CharBuffer[V.VMLI] = V.bankMem.Read(V.ScreenBase + V.VC)
	// fmt.Printf("VMLI: %02X - VC: %02X - Screen Code: %d - Color: %04X\n", V.VMLI, V.VC, V.CharBuffer[V.VMLI], V.ColorBuffer[V.VMLI])
}

func (V *VIC) drawChar(X int, Y int) {
	if V.drawArea && (V.Reg[REG_CTRL1]&DEN > 0) {
		switch V.MODE {
		case 0b00000000:
			V.StandardTextMode(X, Y)
		case 0b00010000: // Multicolor text mode
			V.MulticolTextMode(X, Y)
		case 0b00100000:
			V.StandardBitmapMode(X, Y)
		case 0b00110000:
			V.MulticolBitmapMode(X, Y)
		case 0b01000000: // ECM text mode
			log.Fatal("ECM text mode")
		case 0b01010000: // Invalid text mode
		case 0b01100000: // Invalid bitmap mode 1
		case 0b01110000: // Invalid bitmap mode 2
		}
		V.VMLI++
		V.VC++
	} else if V.visibleArea {
		for column := 0; column < 8; column++ {
			V.graph.DrawPixel(X+column-DisplayOriginX, Y-DisplayOriginY, Colors[V.Reg[REG_BORDER_COL]&0b00001111])
		}
	}
}

func (V *VIC) Run(debug bool) bool {
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
				// fmt.Printf("\nIRQ: %04X - %04X", V.RasterIRQ, uint16(V.BeamY))
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
		if !V.BA {
			V.readVideoMatrix()
		}
	case 16: // Debut de la zone d'affichage
		V.drawChar(V.BeamX, V.BeamY)
		if !V.BA {
			V.readVideoMatrix()
		}
	case 17:
		V.drawChar(V.BeamX, V.BeamY)
		if !V.BA {
			V.readVideoMatrix()
		}
	case 18:
		V.drawChar(V.BeamX, V.BeamY)
		if !V.BA {
			V.readVideoMatrix()
		}
	case 19:
		V.drawChar(V.BeamX, V.BeamY)
		if !V.BA {
			V.readVideoMatrix()
		}
	case 20:
		V.drawChar(V.BeamX, V.BeamY)
		if !V.BA {
			V.readVideoMatrix()
		}
	case 21:
		V.drawChar(V.BeamX, V.BeamY)
		if !V.BA {
			V.readVideoMatrix()
		}
	case 22:
		V.drawChar(V.BeamX, V.BeamY)
		if !V.BA {
			V.readVideoMatrix()
		}
	case 23:
		V.drawChar(V.BeamX, V.BeamY)
		if !V.BA {
			V.readVideoMatrix()
		}
	case 24:
		V.drawChar(V.BeamX, V.BeamY)
		if !V.BA {
			V.readVideoMatrix()
		}
	case 25:
		V.drawChar(V.BeamX, V.BeamY)
		if !V.BA {
			V.readVideoMatrix()
		}
	case 26:
		V.drawChar(V.BeamX, V.BeamY)
		if !V.BA {
			V.readVideoMatrix()
		}
	case 27:
		V.drawChar(V.BeamX, V.BeamY)
		if !V.BA {
			V.readVideoMatrix()
		}
	case 28:
		V.drawChar(V.BeamX, V.BeamY)
		if !V.BA {
			V.readVideoMatrix()
		}
	case 29:
		V.drawChar(V.BeamX, V.BeamY)
		if !V.BA {
			V.readVideoMatrix()
		}
	case 30:
		V.drawChar(V.BeamX, V.BeamY)
		if !V.BA {
			V.readVideoMatrix()
		}
	case 31:
		V.drawChar(V.BeamX, V.BeamY)
		if !V.BA {
			V.readVideoMatrix()
		}
	case 32:
		V.drawChar(V.BeamX, V.BeamY)
		if !V.BA {
			V.readVideoMatrix()
		}
	case 33:
		V.drawChar(V.BeamX, V.BeamY)
		if !V.BA {
			V.readVideoMatrix()
		}
	case 34:
		V.drawChar(V.BeamX, V.BeamY)
		if !V.BA {
			V.readVideoMatrix()
		}
	case 35:
		V.drawChar(V.BeamX, V.BeamY)
		if !V.BA {
			V.readVideoMatrix()
		}
	case 36:
		V.drawChar(V.BeamX, V.BeamY)
		if !V.BA {
			V.readVideoMatrix()
		}
	case 37:
		V.drawChar(V.BeamX, V.BeamY)
		if !V.BA {
			V.readVideoMatrix()
		}
	case 38:
		V.drawChar(V.BeamX, V.BeamY)
		if !V.BA {
			V.readVideoMatrix()
		}
	case 39:
		V.drawChar(V.BeamX, V.BeamY)
		if !V.BA {
			V.readVideoMatrix()
		}
	case 40:
		V.drawChar(V.BeamX, V.BeamY)
		if !V.BA {
			V.readVideoMatrix()
		}
	case 41:
		V.drawChar(V.BeamX, V.BeamY)
		if !V.BA {
			V.readVideoMatrix()
		}
	case 42:
		V.drawChar(V.BeamX, V.BeamY)
		if !V.BA {
			V.readVideoMatrix()
		}
	case 43:
		V.drawChar(V.BeamX, V.BeamY)
		if !V.BA {
			V.readVideoMatrix()
		}
	case 44:
		V.drawChar(V.BeamX, V.BeamY)
		if !V.BA {
			V.readVideoMatrix()
		}
	case 45:
		V.drawChar(V.BeamX, V.BeamY)
		if !V.BA {
			V.readVideoMatrix()
		}
	case 46:
		V.drawChar(V.BeamX, V.BeamY)
		if !V.BA {
			V.readVideoMatrix()
		}
	case 47:
		V.drawChar(V.BeamX, V.BeamY)
		if !V.BA {
			V.readVideoMatrix()
		}
	case 48:
		V.drawChar(V.BeamX, V.BeamY)
		if !V.BA {
			V.readVideoMatrix()
		}
	case 49:
		V.drawChar(V.BeamX, V.BeamY)
		if !V.BA {
			V.readVideoMatrix()
		}
	case 50:
		V.drawChar(V.BeamX, V.BeamY)
		if !V.BA {
			V.readVideoMatrix()
		}
	case 51:
		V.drawChar(V.BeamX, V.BeamY)
		if !V.BA {
			V.readVideoMatrix()
		}
	case 52:
		V.drawChar(V.BeamX, V.BeamY)
		if !V.BA {
			V.readVideoMatrix()
		}
	case 53:
		V.drawChar(V.BeamX, V.BeamY)
		if !V.BA {
			V.readVideoMatrix()
		}
	case 54: // Dernier lecture de la matrice video ram
		V.drawChar(V.BeamX, V.BeamY)
		if !V.BA {
			V.readVideoMatrix()
		}
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
			V.graph.UpdateFrame()
		}
	}
	if debug {
		V.graph.UpdateFrame()
	}
	return V.BA
}

func (V *VIC) Dump(addr uint16) {
	fmt.Printf("Bank: %d - VideoBase: %04X - CharBase: %04X", V.BankSel, V.ScreenBase, V.CharBase)
	V.bankMem.Show()
	V.bankMem.Dump(addr)
}

func (V *VIC) Stats() {
	banks := [4]uint16{BankStart0, BankStart1, BankStart2, BankStart3}

	fmt.Printf("VIC:\n")
	fmt.Printf("Bank: %d - VideoBase: %04X (%04X) - CharBase: %04X (%04X)\n", V.BankSel, V.ScreenBase, banks[V.BankSel]+V.ScreenBase, V.CharBase, banks[V.BankSel]+V.CharBase)
	fmt.Printf("RstX: %04X - RstY: %04X - RC: %02d - VC: %03X - VCBase: %03X - VMLI: %02d\n", V.BeamX, V.BeamY, V.RC, V.VC, V.VCBASE, V.VMLI)
	fmt.Printf("IRQ Line: ")
	if V.Reg[REG_IRQ]&0b10000000 > 0 {
		fmt.Printf("On")
	} else {
		fmt.Printf("Off")
	}
	fmt.Printf(" - IRQ Enabled: ")
	if V.Reg[REG_IRQ_ENABLED]&0b00001111 > 0 {
		fmt.Printf("Yes")
	} else {
		fmt.Printf("None")
	}
	fmt.Printf(" - Raster IRQ: %04X\n", V.RasterIRQ)
}
