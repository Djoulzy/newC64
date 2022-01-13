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
	if elapsed > time.Microsecond {
		log.Printf("%s took %s", name, elapsed)
		C.Disassemble()
		fmt.Printf("\n")
	}
}

func (C *CPU) Reset() {
	C.A = 0xAA
	C.X = 0
	C.Y = 0
	C.S = 0b00100000
	C.SP = 0xFF

	C.IRQ = 0
	C.NMI = 0

	C.ram.Clear(pla906114.RAM)
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

func (C *CPU) irq() {
	//fmt.Printf("\nInterrupt ... Raster: %04X", C.readRasterLine())
	// C.IRQ = 0
	C.pushWordStack(C.PC)
	C.pushByteStack(C.S)
	C.setI(true)
	C.PC = C.readWord(0xFFFE)
}

func (C *CPU) nmi() {
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
	if C.cycleCount == C.inst.cycles {
		// defer C.timeTrack(time.Now(), "ComputeInstruction")
		C.State = ReadInstruction
		C.inst.action()
		if C.NMI > 0 {
			// log.Printf("NMI")
			C.nmi()
		}
		if (C.IRQ > 0) && (C.S & ^I_mask) == 0 {
			// log.Printf("IRQ")
			C.irq()
		}
	}
}

func (C *CPU) NextCycle() {
	var ok bool

	C.cycleCount++
	switch C.State {
	case Idle:
		C.cycleCount = 0
		C.State++
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
		if C.inst.bytes > 1 {
			C.State = ReadOperLO
		} else {
			C.State = Compute
			C.PC += 1
			C.ComputeInstruction()
		}
	case ReadOperLO:
		C.oper = uint16(C.ram.Read(C.PC + 1))
		if C.conf.Disassamble {
			C.instDump += fmt.Sprintf(" %02X", C.ram.Read(C.PC+1))
		}
		if C.inst.bytes > 2 {
			C.State = ReadOperHI
		} else {
			C.State = Compute
			C.PC += 2
			C.ComputeInstruction()
		}
	case ReadOperHI:
		if C.conf.Disassamble {
			C.instDump += fmt.Sprintf(" %02X", C.ram.Read(C.PC+2))
		}
		C.oper += uint16(C.ram.Read(C.PC+2)) << 8
		C.State = Compute
		C.PC += 3
		C.ComputeInstruction()
	case Compute:
		C.ComputeInstruction()
	default:
		log.Fatal("Unknown CPU state\n")
	}
}
