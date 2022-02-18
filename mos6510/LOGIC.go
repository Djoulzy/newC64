package mos6510

import (
	"fmt"
	"log"
)

func (C *CPU) and() {
	var crossed bool

	switch C.Inst.addr {
	case immediate:
		C.A &= byte(C.oper)
	case zeropage:
		C.A &= C.ram.Read(C.oper)
	case zeropageX:
		C.A &= C.ram.Read(C.oper + uint16(C.X))
	case absolute:
		C.A &= C.ram.Read(C.oper)
	case absoluteX:
		C.cross_oper = C.oper + uint16(C.X)
		if C.oper&0xFF00 == C.cross_oper&0xFF00 {
			C.A &= C.ram.Read(C.cross_oper)
		} else {
			C.Inst.addr = CrossPage
			C.State = Compute
			C.Inst.Cycles++
			return
		}
	case absoluteY:
		C.cross_oper = C.oper + uint16(C.Y)
		if C.oper&0xFF00 == C.cross_oper&0xFF00 {
			C.A &= C.ram.Read(C.cross_oper)
		} else {
			C.Inst.addr = CrossPage
			C.State = Compute
			C.Inst.Cycles++
			return
		}
	case indirectX:
		C.A &= C.ReadIndirectX(C.oper)
	case indirectY:
		C.cross_oper = C.GetIndirectYAddr(C.oper, &crossed)
		if crossed {
			C.A &= C.ram.Read(C.cross_oper)
		} else {
			C.Inst.addr = CrossPage
			C.State = Compute
			C.Inst.Cycles++
			return
		}
	case CrossPage:
		C.A &= C.ram.Read(C.cross_oper)
	default:
		log.Fatal("Bad addressing mode")
	}
	C.updateN(C.A)
	C.updateZ(C.A)

}

func (C *CPU) asl() {
	var val byte
	var oper byte
	var dest uint16

	switch C.Inst.addr {
	case implied:
		oper = C.A
		val = oper << 1
		C.A = val
	case zeropage:
		oper = C.ram.Read(C.oper)
		C.ram.Write(C.oper, oper)
		val = oper << 1
		C.ram.Write(C.oper, val)
	case zeropageX:
		dest = C.oper + uint16(C.X)
		oper = C.ram.Read(dest)
		C.ram.Write(dest, oper)
		val = oper << 1
		C.ram.Write(dest, val)
	case absolute:
		oper = C.ram.Read(C.oper)
		C.ram.Write(C.oper, oper)
		val = oper << 1
		C.ram.Write(C.oper, val)
	case absoluteX:
		C.cross_oper = C.oper + uint16(C.X)
		if C.oper&0xFF00 == C.cross_oper&0xFF00 {
			oper = C.ram.Read(C.cross_oper)
			C.ram.Write(C.cross_oper, oper)
			val = oper << 1
			C.ram.Write(C.cross_oper, val)
		} else {
			C.Inst.addr = CrossPage
			C.State = Compute
			C.Inst.Cycles++
			return
		}
	case CrossPage:
		oper = C.ram.Read(C.cross_oper)
		C.ram.Write(C.cross_oper, oper)
		val = oper << 1
		C.ram.Write(C.cross_oper, val)
	default:
		log.Fatal("Bad addressing mode")
	}
	C.updateN(val)
	C.updateZ(val)
	C.setC(oper&0b10000000 > 1)

}

func (C *CPU) eor() {
	var crossed bool

	switch C.Inst.addr {
	case immediate:
		C.A ^= byte(C.oper)
	case zeropage:
		C.A ^= C.ram.Read(C.oper)
	case zeropageX:
		C.A ^= C.ram.Read(C.oper + uint16(C.X))
	case absolute:
		C.A ^= C.ram.Read(C.oper)
	case absoluteX:
		C.cross_oper = C.oper + uint16(C.X)
		if C.oper&0xFF00 == C.cross_oper&0xFF00 {
			C.A ^= C.ram.Read(C.cross_oper)
		} else {
			C.Inst.addr = CrossPage
			C.State = Compute
			C.Inst.Cycles++
			return
		}
	case absoluteY:
		C.cross_oper = C.oper + uint16(C.Y)
		if C.oper&0xFF00 == C.cross_oper&0xFF00 {
			C.A ^= C.ram.Read(C.cross_oper)
		} else {
			C.Inst.addr = CrossPage
			C.State = Compute
			C.Inst.Cycles++
			return
		}
	case indirectX:
		C.A ^= C.ReadIndirectX(C.oper)
	case indirectY:
		C.cross_oper = C.GetIndirectYAddr(C.oper, &crossed)
		if crossed {
			C.A ^= C.ram.Read(C.cross_oper)
		} else {
			C.Inst.addr = CrossPage
			C.State = Compute
			C.Inst.Cycles++
			return
		}
	case CrossPage:
		C.A ^= C.ram.Read(C.cross_oper)
	default:
		log.Fatal("Bad addressing mode")
	}
	C.updateN(C.A)
	C.updateZ(C.A)

}

func (C *CPU) lsr() {
	var val byte
	var oper byte
	var dest uint16

	switch C.Inst.addr {
	case implied:
		C.setC(C.A&0x01 == 0x01)
		val = C.A >> 1
		C.A = val
	case zeropage:
		val = C.ram.Read(C.oper)
		C.ram.Write(C.oper, val)
		C.setC(val&0x01 == 0x01)
		val >>= 1
		C.ram.Write(C.oper, val)
	case zeropageX:
		dest = C.oper + uint16(C.X)
		val = C.ram.Read(dest)
		C.ram.Write(dest, val)
		C.setC(val&0x01 == 0x01)
		val >>= 1
		C.ram.Write(dest, val)
	case absolute:
		val = C.ram.Read(C.oper)
		C.ram.Write(C.oper, val)
		C.setC(val&0x01 == 0x01)
		val >>= 1
		C.ram.Write(C.oper, val)
	case absoluteX:
		C.cross_oper = C.oper + uint16(C.X)
		if C.oper&0xFF00 == C.cross_oper&0xFF00 {
			oper = C.ram.Read(C.cross_oper)
			C.ram.Write(C.cross_oper, oper)
			val = oper >> 1
			C.ram.Write(C.cross_oper, val)
		} else {
			C.Inst.addr = CrossPage
			C.State = Compute
			C.Inst.Cycles++
			return
		}
	case CrossPage:
		oper = C.ram.Read(C.cross_oper)
		C.ram.Write(C.cross_oper, oper)
		val = oper >> 1
		C.ram.Write(C.cross_oper, val)
	default:
		log.Fatal("Bad addressing mode")
	}
	C.setN(false)
	C.updateZ(byte(val))
}

func (C *CPU) ora() {
	switch C.Inst.addr {
	case immediate:
		C.A |= byte(C.oper)
	case zeropage:
		C.A |= C.ram.Read(C.oper)
	case zeropageX:
		C.A |= C.ram.Read(C.oper + uint16(C.X))
	case absolute:
		C.A |= C.ram.Read(C.oper)
	case absoluteX:
		C.A |= C.ram.Read(C.oper + uint16(C.X))
	case absoluteY:
		C.A |= C.ram.Read(C.oper + uint16(C.Y))
	case indirectX:
		C.A |= C.ReadIndirectX(C.oper)
	case indirectY:
		C.A |= C.ReadIndirectY(C.oper)
	default:
		log.Fatal("Bad addressing mode")
	}
	C.updateN(C.A)
	C.updateZ(C.A)

}

func (C *CPU) rla() {
	fmt.Printf("%s\nNot implemented: %v\n", C.Disassemble(), C.Inst)
}

func (C *CPU) rol() {
	var val uint16
	var dest uint16

	switch C.Inst.addr {
	case implied:
		val = uint16(C.A) << 1
		if C.issetC() {
			val++
		}
		C.A = byte(val)
	case zeropage:
		val = uint16(C.ram.Read(C.oper))
		C.ram.Write(C.oper, byte(val))
		val <<= 1
		if C.issetC() {
			val++
		}
		C.ram.Write(C.oper, byte(val))
	case zeropageX:
		dest = C.oper + uint16(C.X)
		val = uint16(C.ram.Read(dest))
		C.ram.Write(dest, byte(val))
		val <<= 1
		if C.issetC() {
			val++
		}
		C.ram.Write(dest, byte(val))
	case absolute:
		val = uint16(C.ram.Read(C.oper))
		C.ram.Write(C.oper, byte(val))
		val <<= 1
		if C.issetC() {
			val++
		}
		C.ram.Write(C.oper, byte(val))
	case absoluteX:
		dest = C.oper + uint16(C.X)
		val = uint16(C.ram.Read(dest))
		C.ram.Write(dest, byte(val))
		val <<= 1
		if C.issetC() {
			val++
		}
		C.ram.Write(dest, byte(val))
	default:
		log.Fatal("Bad addressing mode")
	}
	C.updateN(byte(val))
	C.updateZ(byte(val))
	C.setC(val > 0x00FF)

}

func (C *CPU) ror() {
	var val byte

	switch C.Inst.addr {
	case implied:
		carry := C.A&0b00000001 > 0
		C.A >>= 1
		if C.issetC() {
			C.A |= 0b10000000
		}
		C.setC(carry)
		val = C.A
	case zeropage:
		val = C.ram.Read(C.oper)
		C.ram.Write(C.oper, val)
		carry := val&0b00000001 > 0
		val >>= 1
		if C.issetC() {
			val |= 0b10000000
		}
		C.setC(carry)
		C.ram.Write(C.oper, val)
	case zeropageX:
		dest := C.oper + uint16(C.X)
		val = C.ram.Read(dest)
		C.ram.Write(dest, val)
		carry := val&0b00000001 > 0
		val >>= 1
		if C.issetC() {
			val |= 0b10000000
		}
		C.setC(carry)
		C.ram.Write(dest, val)
	case absolute:
		val = C.ram.Read(C.oper)
		C.ram.Write(C.oper, val)
		carry := val&0b00000001 > 0
		val >>= 1
		if C.issetC() {
			val |= 0b10000000
		}
		C.setC(carry)
		C.ram.Write(C.oper, val)
	case absoluteX:
		dest := C.oper + uint16(C.X)
		val = C.ram.Read(dest)
		C.ram.Write(dest, val)
		carry := val&0b00000001 > 0
		val >>= 1
		if C.issetC() {
			val |= 0b10000000
		}
		C.setC(carry)
		C.ram.Write(dest, val)
	default:
		log.Fatal("Bad addressing mode")
	}
	C.updateN(byte(val))
	C.updateZ(byte(val))

}

func (C *CPU) sax() {
	fmt.Printf("Not implemented: %v\n", C.Inst)
}

func (C *CPU) slo() {
	fmt.Printf("Not implemented: %v\n", C.Inst)
}

func (C *CPU) sre() {
	fmt.Printf("Not implemented: %v\n", C.Inst)
}
