package mos6510

import (
	"fmt"
	"log"
)

// switch C.inst.addr {
// case implied:
// case immediate:
// case relative:
// case zeropage:
// case zeropageX:
// case zeropageY:
// case absolute:
// case absoluteX:
// case absoluteY:
// case indirect:
// case indirectX:
// case indirectY:
// default:
// 	log.Fatal("Bad addressing mode")
// }

func (C *CPU) lda() {
	switch C.inst.addr {
	case immediate:
		C.A = byte(C.oper)
	case zeropage:
		C.A = C.ram.Read(C.oper)
	case zeropageX:
		C.A = C.ram.Read(C.oper + uint16(C.X))
	case absolute:
		C.A = C.ram.Read(C.oper)
	case absoluteX:
		C.A = C.ram.Read(C.oper + uint16(C.X))
	case absoluteY:
		C.A = C.ram.Read(C.oper + uint16(C.Y))
	case indirectX:
		C.A = C.ReadIndirectX(C.oper)
	case indirectY:
		C.A = C.ReadIndirectY(C.oper)
	default:
		log.Fatal("Bad addressing mode")
	}
	C.updateN(C.A)
	C.updateZ(C.A)
	fmt.Printf("\n")
}

func (C *CPU) sta() {
	switch C.inst.addr {
	case zeropage:
		C.ram.Write(C.oper, C.A)
	case zeropageX:
		C.ram.Write(C.oper+uint16(C.X), C.A)
	case absolute:
		C.ram.Write(C.oper, C.A)
	case absoluteX:
		C.ram.Write(C.oper+uint16(C.X), C.A)
	case absoluteY:
		C.ram.Write(C.oper+uint16(C.Y), C.A)
	case indirectX:
		C.WriteIndirectX(C.oper, C.A)
	case indirectY:
		C.WriteIndirectY(C.oper, C.A)
	default:
		log.Fatal("Bad addressing mode")
	}
	fmt.Printf("\n")
}

func (C *CPU) ldx() {
	switch C.inst.addr {
	case immediate:
		C.X = byte(C.oper)
	case zeropage:
		C.X = C.ram.Read(C.oper)
	case zeropageY:
		C.X = C.ram.Read(C.oper + uint16(C.Y))
	case absolute:
		C.X = C.ram.Read(C.oper)
	case absoluteY:
		C.X = C.ram.Read(C.oper + uint16(C.Y))
	default:
		log.Fatal("Bad addressing mode")
	}
	C.updateN(C.X)
	C.updateZ(C.X)
	fmt.Printf("\n")
}

func (C *CPU) stx() {
	switch C.inst.addr {
	case zeropage:
		C.ram.Write(C.oper, C.X)
	case zeropageY:
		C.ram.Write(C.oper+uint16(C.Y), C.X)
	case absolute:
		C.ram.Write(C.oper, C.X)
	default:
		log.Fatal("Bad addressing mode")
	}
	fmt.Printf("\n")
}

func (C *CPU) ldy() {
	switch C.inst.addr {
	case immediate:
		C.Y = byte(C.oper)
	case zeropage:
		C.Y = C.ram.Read(C.oper)
	case zeropageX:
		C.Y = C.ram.Read(C.oper + uint16(C.X))
	case absolute:
		C.Y = C.ram.Read(C.oper)
	case absoluteX:
		C.Y = C.ram.Read(C.oper + uint16(C.X))
	default:
		log.Fatal("Bad addressing mode")
	}
	C.updateN(C.Y)
	C.updateZ(C.Y)
	fmt.Printf("\n")
}

func (C *CPU) sty() {
	switch C.inst.addr {
	case zeropage:
		C.ram.Write(C.oper, C.Y)
	case zeropageX:
		C.ram.Write(C.oper+uint16(C.X), C.Y)
	case absolute:
		C.ram.Write(C.oper, C.Y)
	default:
		log.Fatal("Bad addressing mode")
	}
	fmt.Printf("\n")
}
