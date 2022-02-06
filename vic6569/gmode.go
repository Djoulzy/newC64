package vic6569

func (V *VIC) StandardTextMode(X int, Y int) {
	var pixelData byte
	var palette [4]byte

	charAddr := (uint16(V.CharBuffer[V.VMLI]) << 3) + uint16(V.RC)
	pixelData = V.bankMem[V.BankSel].Read(V.CharBase + charAddr)
	palette[1] = V.ColorBuffer[V.VMLI] // V.color.Val[V.VC] & 0b00001111
	palette[0] = V.Reg[REG_BGCOLOR_0] & 0b00001111

	for column := 0; column < 8; column++ {
		bit := byte(0b10000000 >> column)
		if pixelData&bit > 0 {
			V.graph.DrawPixel(X+column, Y, Colors[palette[1]])
		} else {
			V.graph.DrawPixel(X+column, Y, Colors[palette[0]])
		}
	}
}

func (V *VIC) MulticolTextMode(X int, Y int) {
	var pixelData byte
	var palette [4]byte

	charAddr := (uint16(V.CharBuffer[V.VMLI]) << 3) + uint16(V.RC)
	pixelData = V.bankMem[V.BankSel].Read(V.CharBase + charAddr)
	MCFlag := V.ColorBuffer[V.VMLI] & 0b00001000

	if MCFlag > 0 {
		palette[0] = V.Reg[REG_BGCOLOR_0] & 0b00001111
		palette[1] = V.Reg[REG_BGCOLOR_1] & 0b00001111
		palette[2] = V.Reg[REG_BGCOLOR_2] & 0b00001111
		palette[3] = V.ColorBuffer[V.VMLI] & 0b00000111

		bit := byte(0b11000000)
		for column := 0; column < 4; column++ {
			colNum := (pixelData & bit) >> byte((3-column)<<1)
			V.graph.DrawPixel(X+(column<<1), Y, Colors[palette[colNum]])
			V.graph.DrawPixel(X+(column<<1)+1, Y, Colors[palette[colNum]])
			bit >>= 2
		}
	} else {
		palette[1] = V.ColorBuffer[V.VMLI] & 0b00000111
		palette[0] = V.Reg[REG_BGCOLOR_0] & 0b00001111

		for column := 0; column < 8; column++ {
			bit := byte(0b10000000 >> column)
			if pixelData&bit > 0 {
				V.graph.DrawPixel(X+column, Y, Colors[palette[1]])
			} else {
				V.graph.DrawPixel(X+column, Y, Colors[palette[0]])
			}
		}
	}
}

func (V *VIC) MulticolBitmapMode(X int, Y int) {
	var pixelData byte
	var palette [4]byte

	pixAddr := (V.CharBase & 0x2000) + (V.VC << 3) + uint16(V.RC)
	colors := V.bankMem[V.BankSel].Read(V.ScreenBase + V.VC)
	palette[0] = V.Reg[REG_BGCOLOR_0] & 0b00001111
	palette[1] = colors >> 4
	palette[2] = colors & 0b00001111
	palette[3] = V.ColorBuffer[V.VMLI] // V.color.Val[V.VC] & 0b00001111
	pixelData = V.bankMem[V.BankSel].Read(pixAddr)

	bit := byte(0b11000000)
	for column := 0; column < 4; column++ {
		colNum := (pixelData & bit) >> byte((3-column)<<1)
		V.graph.DrawPixel(X+(column<<1), Y, Colors[palette[colNum]])
		V.graph.DrawPixel(X+(column<<1)+1, Y, Colors[palette[colNum]])
		bit >>= 2
	}
}

func (V *VIC) StandardBitmapMode(X int, Y int) {
	var pixelData byte
	var palette [4]byte

	pixAddr := (V.CharBase & 0x2000) + (V.VC << 3) + uint16(V.RC)
	colors := V.bankMem[V.BankSel].Read(V.ScreenBase + V.VC)
	palette[1] = colors >> 4
	palette[0] = colors & 0b00001111
	pixelData = V.bankMem[V.BankSel].Read(pixAddr)
	// log.Printf("Read VM %04X", pixAddr)

	for column := 0; column < 8; column++ {
		bit := byte(0b10000000 >> column)
		if pixelData&bit > 0 {
			V.graph.DrawPixel(X+column, Y, Colors[palette[1]])
		} else {
			V.graph.DrawPixel(X+column, Y, Colors[palette[0]])
		}
	}
}
