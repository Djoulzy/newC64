package mos6510

import (
	"fmt"
	"log"
	"newC64/confload"
	"newC64/pla906114"
	"time"
)

func (C *CPU) timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	// if elapsed > time.Microsecond {
		log.Printf("Phase %d - %s - %s took %s", C.State, C.inst.name, name, elapsed)
	// }
}

func (C *CPU) Reset() {
	C.A = 0xAA
	C.X = 0
	C.Y = 0
	C.S = 0b00100000
	C.SP = 0xFF

	C.IRQ_pin = 0
	C.NMI_pin = 0

	// PLA Settings (Bank switching)
	// C.ram.Write(0x0000, 0x2F)
	C.ram.Write(0x0001, 0x1F)

	C.State = ReadInstruction
	// Cold Start:
	C.PC = C.readWord(COLDSTART_Vector)
	fmt.Printf("mos6510 - PC: %04X\n", C.PC)
}

func (C *CPU) Init(mem *pla906114.PLA, clock *uint16, conf *confload.ConfigData) {
	fmt.Printf("mos6510 - Init\n")
	C.conf = conf
	C.ClockCycles = clock
	C.ram = mem
	C.stack = C.ram.GetView(StackStart, 256)
	C.initLanguage()
	C.Reset()
}

func (C *CPU) registers() string {
	var i, mask byte
	res := ""
	for i = 0; i < 8; i++ {
		mask = 1 << i
		if C.S&mask == mask {
			res = regString[i] + res
		} else {
			res = "-" + res
		}
	}
	return res
}

func (C *CPU) Disassemble() string {
	var buf, token string

	buf = fmt.Sprintf("%s - A:%c[1;33m%02X%c[0m X:%c[1;33m%02X%c[0m Y:%c[1;33m%02X%c[0m SP:%c[1;33m%02X%c[0m - ", C.registers(), 27, C.A, 27, 27, C.X, 27, 27, C.Y, 27, 27, C.SP, 27)
	buf = fmt.Sprintf("%s%04X: %-8s %03s ", buf, C.InstStart, C.instDump, C.inst.name)
	switch C.inst.addr {
	case implied:
		token = fmt.Sprintf("")
	case immediate:
		token = fmt.Sprintf("#$%02X", C.oper)
	case relative:
		token = fmt.Sprintf("$%02X", C.oper)
	case zeropage:
		token = fmt.Sprintf("$%02X", C.oper)
	case zeropageX:
		token = fmt.Sprintf("$%02X,X", C.oper)
	case zeropageY:
		token = fmt.Sprintf("$%02X,Y", C.oper)
	case absolute:
		token = fmt.Sprintf("$%04X", C.oper)
	case absoluteX:
		token = fmt.Sprintf("$%04X,X", C.oper)
	case absoluteY:
		token = fmt.Sprintf("$%04X,Y", C.oper)
	case indirect:
		token = fmt.Sprintf("($%04X)", C.oper)
	case indirectX:
		token = fmt.Sprintf("($%02X,X)", C.oper)
	case indirectY:
		token = fmt.Sprintf("($%02X),Y", C.oper)
	}
	return fmt.Sprintf("%s%-10s\t", buf, token)
}

//////////////////////////////////
////// Addressage Indirect ///////
//////////////////////////////////

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

//////////////////////////////////
/////// Addressage Relatif ///////
//////////////////////////////////

func (C *CPU) getRelativeAddr(dist uint16) uint16 {
	return uint16(int(C.PC) + int(int8(dist)))
}

//////////////////////////////////
//////////// Read Word ///////////
//////////////////////////////////

func (C *CPU) readWord(addr uint16) uint16 {
	return (uint16(C.ram.Read(addr+1)) << 8) + uint16(C.ram.Read(addr))
}

//////////////////////////////////
//////// Stack Operations ////////
//////////////////////////////////

// Byte
func (C *CPU) pushByteStack(val byte) {
	// if C.SP < 90 {
	// 	os.Exit(1)
	// }
	C.stack[C.SP] = val
	C.SP--
}

func (C *CPU) pullByteStack() byte {
	C.SP++
	// if C.SP == 0x00 {
	// 	C.ram.DumpStack(C.SP)
	// 	log.Fatal("Stack overflow")
	// }
	return C.stack[C.SP]
}

// Word
func (C *CPU) pushWordStack(val uint16) {
	low := byte(val)
	hi := byte(val >> 8)
	C.pushByteStack(hi)
	C.pushByteStack(low)
}

func (C *CPU) pullWordStack() uint16 {
	low := C.pullByteStack()
	hi := uint16(C.pullByteStack()) << 8
	return hi + uint16(low)
}

//////////////////////////////////
/////////// Interrupts ///////////
//////////////////////////////////

func (C *CPU) IRQ() {
	//fmt.Printf("\nInterrupt ... Raster: %04X", C.readRasterLine())
	// C.IRQ = 0
	C.pushWordStack(C.PC)
	C.pushByteStack(C.S)
	C.setI(true)
	C.PC = C.readWord(0xFFFE)
}

func (C *CPU) NMI() {
	//fmt.Printf("\nInterrupt ... Raster: %04X", C.readRasterLine())
	// C.NMI = 0
	C.pushWordStack(C.PC)
	C.pushByteStack(C.S)
	C.PC = C.readWord(0xFFFA)
}

//////////////////////////////////
///////////// Running ////////////
//////////////////////////////////

func (C *CPU) GoTo(addr uint16) {
	C.PC = addr
}

func (C *CPU) ComputeInstruction() {
	// if C.cycleCount == C.inst.cycles {
	// defer C.timeTrack(time.Now(), "ComputeInstruction")
	C.State = ReadInstruction
	C.inst.action()
	if C.cycleCount != C.inst.cycles {
		log.Printf("%s - Wanted: %d - Getting: %d\n", C.Disassemble(), C.inst.cycles, C.cycleCount)
	}
	if C.cycleCount == C.inst.cycles {
		C.State = ReadInstruction
	}
	// }
}

func (C *CPU) NextCycle() {
	defer C.timeTrack(time.Now(), "ComputeInstruction")
	var ok bool

	C.cycleCount++
	// fmt.Printf("%d - %d\n", C.cycleCount, C.State)
	switch C.State {
	case Idle:
		C.cycleCount = 0
		C.State++

	////////////////////////////////////////////////
	// Cycle 1
	////////////////////////////////////////////////
	case ReadInstruction:
		C.cycleCount = 1
		C.InstStart = C.PC
		if C.conf.Disassamble {
			C.instDump = fmt.Sprintf("%02X", C.ram.Read(C.PC))
		}
		if C.inst, ok = mnemonic[C.ram.Read(C.PC)]; !ok {
			log.Printf(fmt.Sprintf("Unknown instruction: %02X at %04X\n", C.ram.Read(C.PC), C.PC))
			// C.State = Idle
		}
		if C.inst.addr == implied {
			C.State = Compute
			C.PC += 1
		} else {
			C.State = ReadOperLO
		}

	////////////////////////////////////////////////
	// Cycle 2
	////////////////////////////////////////////////
	case ReadOperLO:
		C.oper = uint16(C.ram.Read(C.PC + 1))
		if C.conf.Disassamble {
			C.instDump += fmt.Sprintf(" %02X", C.ram.Read(C.PC+1))
		}
		switch C.inst.addr {
		case relative:
			fallthrough
		case immediate:
			C.State = Compute
			C.PC += 2
			if C.inst.cycles == 2 {
				C.ComputeInstruction()
			}
		case absolute:
			fallthrough
		case indirect:
			fallthrough
		case absoluteX:
			fallthrough
		case absoluteY:
			C.State = ReadOperHI
		case zeropage:
			fallthrough
		case zeropageX:
			fallthrough
		case zeropageY:
			fallthrough
		case indirectX:
			fallthrough
		case indirectY:
			C.State = ReadZP
		default:
			log.Fatal("Erreur de cycle")
		}

	////////////////////////////////////////////////
	// Cycle 3
	////////////////////////////////////////////////
	case ReadZP:
		C.PC += 2
		switch C.inst.addr {
		case zeropage:
			C.State = Compute
			if C.inst.cycles == 3 {
				C.ComputeInstruction()
			}
		case zeropageX:
			fallthrough
		case zeropageY:
			C.State = ReadZP_XY
		case indirectX:
			fallthrough
		case indirectY:
			C.State = ReadIndXY_LO
		default:
			log.Fatal("Erreur de cycle")
		}

	case ReadOperHI: // Cycle 3
		C.oper += uint16(C.ram.Read(C.PC+2)) << 8
		C.PC += 3
		if C.conf.Disassamble {
			C.instDump += fmt.Sprintf(" %02X", C.ram.Read(C.PC+2))
		}
		switch C.inst.addr {
		case absolute:
			C.State = Compute
			if C.inst.cycles == 3 {
				C.ComputeInstruction()
			}
		case absoluteX:
			fallthrough
		case absoluteY:
			C.State = ReadAbsXY
		case indirect:
			C.State = ReadIndirect
		default:
			log.Fatal("Erreur de cycle")
		}

	////////////////////////////////////////////////
	// Cycle 4
	////////////////////////////////////////////////
	case ReadZP_XY: // Cycle 4
		switch C.inst.addr {
		case zeropageX:
			fallthrough
		case zeropageY:
			C.State = Compute
			if C.inst.cycles == 4 {
				C.ComputeInstruction()
			}
		default:
			log.Fatal("Erreur de cycle")
		}

	case ReadIndXY_LO: // Cycle 4
		switch C.inst.addr {
		case indirectX:
			C.State = ReadIndXY_HI
		case indirectY:
			C.State = ReadIndXY_HI
		default:
			log.Fatal("Erreur de cycle")
		}

	case ReadIndirect: // Cycle 4
		C.State = Compute

	case ReadAbsXY: // Cycle 4
		switch C.inst.addr {
		case absoluteX:
			fallthrough
		case absoluteY:
			C.State = Compute
			if C.inst.cycles == 4 {
				C.ComputeInstruction()
			}
		default:
			log.Fatal("Erreur de cycle")
		}

	////////////////////////////////////////////////
	// Cycle 5
	////////////////////////////////////////////////

	case ReadIndXY_HI:
		switch C.inst.addr {
		case indirectX:
			C.State = Compute
		case indirectY:
			C.State = Compute
			if C.inst.cycles == 5 {
				C.ComputeInstruction()
			}
		default:
			log.Fatal("Erreur de cycle")
		}

	////////////////////////////////////////////////
	// Exec
	////////////////////////////////////////////////
	case Compute:
		// if C.cycleCount > C.inst.cycles {
		// 	log.Printf("%s - Wanted: %d - Getting: %d\n", C.Disassemble(), C.inst.cycles, C.cycleCount)
		// }
		if C.inst.cycles == C.cycleCount {
			C.ComputeInstruction()
		}
	default:
		log.Fatal("Unknown CPU state\n")
	}
}
