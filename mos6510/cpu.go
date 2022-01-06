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
	// Cold Start:
	// C.PC = C.readWord(COLDSTART_Vector)
	C.PC = 0xE000
	fmt.Printf("mos6510 - PC: %04X\n", C.PC)
}

func (C *CPU) Init(mem interface{}) {
	fmt.Printf("mos6510 - Init\n")
	C.ram = mem.(memory)
	C.ram.Init()
	C.stack = (C.ram.GetView(StackStart, 256)).(memory)
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
	C.stack.Write(uint16(C.SP), val)
	C.SP--
}

func (C *CPU) pullByteStack() byte {
	if C.SP == 0xFF {
		log.Fatal("Stack overflow")
	}
	C.SP++
	return C.stack.Read(uint16(C.SP))
}

// Word
func (C *CPU) pushWordStack(val uint16) {
	C.pushByteStack(byte(val >> 8)) // HI
	C.pushByteStack(byte(val))      // LO
}

func (C *CPU) pullWordStack() uint16 {
	low := C.pullByteStack()
	hi := uint16(C.pullByteStack()) << 8
	return hi + uint16(low)
}

//////////////////////////////////
///////////// Running ////////////
//////////////////////////////////

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
