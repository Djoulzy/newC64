package mos6510

import "fmt"

func (C *CPU) Reset() {
	C.A = 0xAA
	C.X = 0
	C.Y = 0
	C.S = 0b00100000

	C.PC = 0xFF00
	C.SP = 0xFF
}

func (C *CPU) Init() {
	fmt.Printf("mos6510 - Init\n")
}

func (C *CPU) NextCycle() {

}

// var output = ""

// func (C *CPU) readRasterLine() uint16 {
// 	val := uint16(C.ram.Mem[0xD011].Zone[mem.IO]&0b10000000) << 8
// 	val += uint16(C.ram.Mem[0xD012].Zone[mem.IO])
// 	return val
// }

// func (C *CPU) timeTrack(start time.Time) {
// 	elapsed := time.Since(start)
// 	log.Printf("%s", elapsed)
// }

// //////////////////////////////////
// //////// Stack Operations ////////
// //////////////////////////////////

// // Word
// func (C *CPU) pushWordStack(val uint16) {
// 	low := byte(val)
// 	hi := byte(val >> 8)
// 	C.pushByteStack(hi)
// 	C.pushByteStack(low)
// }

// func (C *CPU) pullWordStack() uint16 {
// 	low := C.pullByteStack()
// 	hi := uint16(C.pullByteStack()) << 8
// 	return hi + uint16(low)
// }

// // Byte
// func (C *CPU) pushByteStack(val byte) {
// 	C.ram.Stack[C.SP].Zone[mem.RAM] = val
// 	C.SP--

// }

// func (C *CPU) pullByteStack() byte {
// 	C.SP++
// 	// if C.SP > 0xFF {
// 	// 	panic("Stack overflow")
// 	// }

// 	return C.ram.Stack[C.SP].Zone[mem.RAM]
// }

// //////////////////////////////////
// ////// Addressage Indirect ///////
// //////////////////////////////////

// // https://stackoverflow.com/questions/46262435/indirect-y-indexed-addressing-mode-in-mos-6502
// // http://www.emulator101.com/6502-addressing-modes.html

// func (C *CPU) Indirect_index_Y(addr byte, y byte) uint16 {
// 	wordZP := C.readWord(uint16(addr)) + uint16(y)
// 	return wordZP
// }

// func (C *CPU) Indexed_indirect_X(addr byte, x byte) uint16 {
// 	wordZP := C.readWord(uint16(addr + x))
// 	return wordZP
// }

// //////////////////////////////////
// /////// Memory Operations ////////
// //////////////////////////////////

// func (C *CPU) readWord(addr uint16) uint16 {
// 	low := C.ram.Read(addr)

// 	value := (uint16(C.ram.Read(addr+1)) << 8) + uint16(low)

// 	return value
// }

// func (C *CPU) readByte(addr uint16) byte {

// 	return C.ram.Read(addr)
// }

// func (C *CPU) writeByte(addr uint16, value byte) {
// 	C.ram.Write(addr, value)

// }

// //////////////////////////////////
// ////////// Read OpCode ///////////
// //////////////////////////////////

// func (C *CPU) fetchWord() uint16 {
// 	low := C.fetchByte()
// 	return (uint16(C.fetchByte()) << 8) + uint16(low)
// }

// func (C *CPU) fetchByte() byte {
// 	value := C.ram.Read(C.PC)
// 	C.PC++
// 	if C.Display {
// 		output = fmt.Sprintf("%s %02X", output, value)
// 	}

// 	return value
// }

// func (C *CPU) exec() {
// 	// if C.exit {
// 	// 	os.Exit(1)
// 	// }
// 	if C.Display {
// 		output = ""
// 		fmt.Printf("\n%08b - A:%c[1;33m%02X%c[0m X:%c[1;33m%02X%c[0m Y:%c[1;33m%02X%c[0m SP:%c[1;33m%02X%c[0m", C.S, 27, C.A, 27, 27, C.X, 27, 27, C.Y, 27, 27, C.SP, 27)
// 		fmt.Printf(" RastY: %c[1;31m%04X%c[0m RastX: - %c[1;31m%04X%c[0m:", 27, C.readRasterLine(), 27, 27, C.PC, 27)
// 	}
// 	Mnemonic[C.fetchByte()]()
// 	// if C.opName == "ToDO" {
// 	// 	fmt.Printf("\n\nToDO : %02X\n\n", opCode)
// 	// 	os.Exit(1)
// 	// }
// 	if C.Display {
// 		fmt.Printf("%c[1;30m%-15s%c[0m %-15s%c[0;32m; (%d) %s%c[0m", 27, output, 27, C.opName, 27, 0, C.debug, 27)
// 		C.debug = ""
// 	}
// }

// func (C *CPU) SetBreakpoint(bp uint16) {
// 	C.BP = bp
// }

// func (C *CPU) irq() {
// 	//fmt.Printf("\nInterrupt ... Raster: %04X", C.readRasterLine())
// 	// C.IRQ = 0
// 	C.pushWordStack(C.PC)
// 	C.pushByteStack(C.S)
// 	C.setI(true)
// 	C.PC = C.readWord(0xFFFE)
// }

// func (C *CPU) nmi() {
// 	//fmt.Printf("\nInterrupt ... Raster: %04X", C.readRasterLine())
// 	// C.NMI = 0
// 	C.pushWordStack(C.PC)
// 	C.pushByteStack(C.S)
// 	C.PC = C.readWord(0xFFFA)
// }

// func (C *CPU) Init(mem *mem.Memory) {
// 	C.ram = mem
// 	C.BP = 0

// 	C.initLanguage()
// 	C.reset(C.ram)
// 	C.tty, _ = tty.Open()
// 	// Recupere l'addresse de boot du systeme
// 	C.PC = (uint16(C.ram.Read(0xFFFC+1)) << 8) + uint16(C.ram.Read(0xFFFC))
// }

// func (C *CPU) Run() {
// 	// t0 := time.Now()
// 	if C.PC == C.BP {
// 		C.Display = true
// 		C.Step = true
// 	}

// 	C.exec()
// 	if C.NMI > 0 {
// 		// log.Printf("NMI")
// 		C.nmi()
// 	}
// 	if (C.IRQ > 0) && (C.S & ^I_mask) == 0 {
// 		// log.Printf("IRQ")
// 		C.irq()
// 	}

// 	if C.Step {
// 		// C.ram.DumpCIA()
// 	COMMAND:
// 		r, err := C.tty.ReadRune()
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		switch r {
// 		case 'c':
// 			C.mem.DumpCIA()
// 			goto COMMAND
// 		case 'd':
// 			C.mem.Dump(C.Dump, C.Zone)
// 			goto COMMAND
// 		case 's':
// 			fmt.Printf("\n")
// 			C.mem.DumpStack(C.SP, 0)
// 			goto COMMAND
// 		case 'z':
// 			fmt.Printf("\n")
// 			C.mem.Dump(0x0000, mem.RAM)
// 			goto COMMAND
// 		default:
// 		}
// 	}
// }

//////////////////////////////////
//////////// Language ////////////
//////////////////////////////////

// func (C *CPU) CheckMnemonic(code string) {
// 	test := Mnemonic[code]
// }
