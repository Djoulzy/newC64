package mos6510

import (
	"fmt"
	"log"
)

func (C *CPU) dcp() {
	fmt.Printf("Not implemented: %v\n", C.Inst)
}

func (C *CPU) dec() {
	var val byte

	switch C.Inst.addr {
	case zeropage:
		val = C.ram.Read(C.oper) - 1
		C.ram.Write(C.oper, val)
	case zeropageX:
		val = C.ram.Read(C.oper+uint16(C.X)) - 1
		C.ram.Write(C.oper+uint16(C.X), val)
	case absolute:
		val = C.ram.Read(C.oper) - 1
		C.ram.Write(C.oper, val)
	case absoluteX:
		val = C.ram.Read(C.oper+uint16(C.X)) - 1
		C.ram.Write(C.oper+uint16(C.X), val)
	default:
		log.Fatal("Bad addressing mode")
	}
	C.updateN(val)
	C.updateZ(val)

}

func (C *CPU) dex() {
	switch C.Inst.addr {
	case implied:
		C.X -= 1
	default:
		log.Fatal("Bad addressing mode")
	}
	C.updateN(C.X)
	C.updateZ(C.X)

}

func (C *CPU) dey() {
	switch C.Inst.addr {
	case implied:
		C.Y -= 1
	default:
		log.Fatal("Bad addressing mode")
	}
	C.updateN(C.Y)
	C.updateZ(C.Y)

}

func (C *CPU) inc() {
	var val byte

	switch C.Inst.addr {
	case zeropage:
		val = C.ram.Read(C.oper) + 1
		C.ram.Write(C.oper, val)
	case zeropageX:
		val = C.ram.Read(C.oper+uint16(C.X)) + 1
		C.ram.Write(C.oper+uint16(C.X), val)
	case absolute:
		val = C.ram.Read(C.oper) + 1
		C.ram.Write(C.oper, val)
	case absoluteX:
		val = C.ram.Read(C.oper+uint16(C.X)) + 1
		C.ram.Write(C.oper+uint16(C.X), val)
	default:
		log.Fatal("Bad addressing mode")
	}
	C.updateN(val)
	C.updateZ(val)

}

func (C *CPU) inx() {
	switch C.Inst.addr {
	case implied:
		C.X += 1
	default:
		log.Fatal("Bad addressing mode")
	}
	C.updateN(C.X)
	C.updateZ(C.X)

}

func (C *CPU) iny() {
	switch C.Inst.addr {
	case implied:
		C.Y += 1
	default:
		log.Fatal("Bad addressing mode")
	}
	C.updateN(C.Y)
	C.updateZ(C.Y)

}

func (C *CPU) isc() {
	fmt.Printf("%s\nNot implemented: %v\n", C.Disassemble(), C.Inst)
	// log.Fatal("fichier INC_DEC")
}
