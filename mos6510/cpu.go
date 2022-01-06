package mos6510

import (
	"fmt"
	"log"
)

func (C *CPU) Reset() {
	C.A = 0xAA
	C.X = 0
	C.Y = 0
	C.S = 0b00100000
	C.SP = 0xFF

	C.ram.Clear()
	// PLA Settings (Bank switching)
	C.ram.Write(0x0001, 7)

	C.state = readInstruction
	// C.PC = (uint16(C.ram.Read(0xFFFC+1)) << 8) + uint16(C.ram.Read(0xFFFC))
	C.PC = 0xE000
	fmt.Printf("mos6510 - PC: %04X\n", C.PC)
}

func (C *CPU) Init(mem interface{}) {
	fmt.Printf("mos6510 - Init\n")
	C.ram = mem.(memory)
	C.ram.Init()
	C.initLanguage()
	C.Reset()
}

func (C *CPU) disassemble() {
	fmt.Printf("%04X: %03s ", C.instStart, C.inst.name)
	switch C.inst.addr {
	case implied:
		fmt.Printf("\t\t")
	case immediate:
		fmt.Printf("#$%02X\t\t", C.oper)
	case relative:
		fmt.Printf("$%02X\t\t", C.oper)
	case zeropage:
		fmt.Printf("$%02X\t\t", C.oper)
	case zeropageX:
		fmt.Printf("$%02X,X\t\t", C.oper)
	case zeropageY:
		fmt.Printf("$%02X,Y\t\t", C.oper)
	case absolute:
		fmt.Printf("$%04X\t\t", C.oper)
	case absoluteX:
		fmt.Printf("$%04X,X\t", C.oper)
	case absoluteY:
		fmt.Printf("$%04X,Y\t", C.oper)
	case indirect:
		fmt.Printf("($%04X)\t", C.oper)
	case indirectX:
		fmt.Printf("($%02X,X)\t", C.oper)
	case indirectY:
		fmt.Printf("($%02X),Y\t", C.oper)
	}
	fmt.Printf("\t")
}

func (C *CPU) ReadIndirectX(addr uint16) byte {
	dest := addr + uint16(C.X)
	return C.ram.Read((uint16(C.ram.Read(dest+1)) << 8) + uint16(C.ram.Read(dest)))
}

func (C *CPU) ReadIndirectY(addr uint16) byte {
	dest := (uint16(C.ram.Read(addr+1)) << 8) + uint16(C.ram.Read(addr))
	return C.ram.Read(dest + uint16(C.Y))
}

func (C *CPU) WriteIndirectX(addr uint16, val byte) {
	dest := addr + uint16(C.X)
	C.ram.Write((uint16(C.ram.Read(dest+1))<<8)+uint16(C.ram.Read(dest)), val)
}

func (C *CPU) WriteIndirectY(addr uint16, val byte) {
	dest := (uint16(C.ram.Read(addr+1)) << 8) + uint16(C.ram.Read(addr))
	C.ram.Write(dest+uint16(C.Y), val)
}

func (C *CPU) computeInstruction() {
	if C.cycleCount == C.inst.cycles {
		C.state = readInstruction
		C.disassemble()
		C.inst.action()
	}
}

func (C *CPU) NextCycle() {
	var ok bool

	C.cycleCount++
	switch C.state {
	case idle:
		C.cycleCount = 0
		C.state++
	case readInstruction:
		C.cycleCount = 1
		C.instStart = C.PC
		if C.inst, ok = mnemonic[C.ram.Read(C.PC)]; !ok {
			log.Fatal(fmt.Sprintf("Unknown instruction: %02X at %04X\n", C.ram.Read(C.PC), C.PC))
		}
		C.PC++
		if C.inst.bytes > 1 {
			C.state = readOperLO
		} else {
			C.state = compute
			C.computeInstruction()
		}
	case readOperLO:
		C.oper = uint16(C.ram.Read(C.PC))
		C.PC++
		if C.inst.bytes > 2 {
			C.state = readOperHI
		} else {
			C.state = compute
			C.computeInstruction()
		}
	case readOperHI:
		C.oper += uint16(C.ram.Read(C.PC)) << 8
		C.PC++
		C.state = compute
		C.computeInstruction()
	case compute:
		C.computeInstruction()
	default:
		log.Fatal("Unknown CPU state\n")
	}
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
