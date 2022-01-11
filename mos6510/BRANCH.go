package mos6510

import (
	"fmt"
	"log"
)

func (C *CPU) bcc() {
	switch C.inst.addr {
	case relative:
		if !C.issetC() {
			C.PC = C.getRelativeAddr(C.oper)
		}
	default:
		log.Fatal("Bad addressing mode")
	}
}

func (C *CPU) bcs() {
	switch C.inst.addr {
	case relative:
		if C.issetC() {
			C.PC = C.getRelativeAddr(C.oper)
		}
	default:
		log.Fatal("Bad addressing mode")
	}
}

func (C *CPU) beq() {
	switch C.inst.addr {
	case relative:
		if C.issetZ() {
			C.PC = C.getRelativeAddr(C.oper)
		}
	default:
		log.Fatal("Bad addressing mode")
	}
}

func (C *CPU) bmi() {
	switch C.inst.addr {
	case relative:
		if C.issetN() {
			C.PC = C.getRelativeAddr(C.oper)
		}
	default:
		log.Fatal("Bad addressing mode")
	}
}

func (C *CPU) bne() {
	switch C.inst.addr {
	case relative:
		if !C.issetZ() {
			C.PC = C.getRelativeAddr(C.oper)
		}
	default:
		log.Fatal("Bad addressing mode")
	}
}

func (C *CPU) bpl() {
	switch C.inst.addr {
	case relative:
		if !C.issetN() {
			C.PC = C.getRelativeAddr(C.oper)
		}
	default:
		log.Fatal("Bad addressing mode")
	}
}

func (C *CPU) brk() {
	switch C.inst.addr {
	case implied:
		C.pushWordStack(C.PC + 1)
		C.setB(true)
		C.pushByteStack(C.S)
		C.PC = C.readWord(IRQBRK_Vector)
	default:
		log.Fatal("Bad addressing mode")
	}
}

func (C *CPU) bvc() {
	dest := C.getRelativeAddr(C.oper)
	switch C.inst.addr {
	case relative:
		if !C.issetV() {
			C.PC = C.getRelativeAddr(C.oper)
		}
	default:
		log.Fatal("Bad addressing mode")
	}
	if C.conf.Disassamble {
		fmt.Printf("(%04X)\n", dest)
	}
}

func (C *CPU) bvs() {
	dest := C.getRelativeAddr(C.oper)
	switch C.inst.addr {
	case relative:
		if C.issetV() {
			C.PC = C.getRelativeAddr(C.oper)
		}
	default:
		log.Fatal("Bad addressing mode")
	}
	if C.conf.Disassamble {
		fmt.Printf("(%04X)\n", dest)
	}
}

func (C *CPU) jmp() {
	switch C.inst.addr {
	case absolute:
		C.PC = C.oper
	case indirect:
		C.PC = C.readWord(C.oper)
	default:
		log.Fatal("Bad addressing mode")
	}

}

func (C *CPU) jsr() {
	switch C.inst.addr {
	case absolute:
		C.pushWordStack(C.InstStart + 2)
		C.PC = C.oper
	default:
		log.Fatal("Bad addressing mode")
	}

}

func (C *CPU) rti() {
	switch C.inst.addr {
	case implied:
		C.S = C.pullByteStack()
		C.setB(false)
		C.setU(false)
		C.PC = C.pullWordStack()
	default:
		log.Fatal("Bad addressing mode")
	}

}

func (C *CPU) rts() {
	switch C.inst.addr {
	case implied:
		C.PC = C.pullWordStack() + 1
	default:
		log.Fatal("Bad addressing mode")
	}

}

func (C *CPU) nop() {
	switch C.inst.addr {
	case implied:
		fallthrough
	case immediate:
		fallthrough
	case zeropage:
		fallthrough
	case zeropageX:
		fallthrough
	case absolute:
		fallthrough
	case absoluteX:
	default:
		log.Fatal("Bad addressing mode")
	}

}
