package mos6510

import (
	"log"
)

func (C *CPU) tax() {
	switch C.Inst.addr {
	case implied:
		C.X = C.A
	default:
		log.Fatal("Bad addressing mode")
	}
	C.updateN(C.X)
	C.updateZ(C.X)

}

func (C *CPU) tay() {
	switch C.Inst.addr {
	case implied:
		C.Y = C.A
	default:
		log.Fatal("Bad addressing mode")
	}
	C.updateN(C.Y)
	C.updateZ(C.Y)

}

func (C *CPU) tsx() {
	switch C.Inst.addr {
	case implied:
		C.X = C.SP
	default:
		log.Fatal("Bad addressing mode")
	}
	C.updateN(C.X)
	C.updateZ(C.X)

}

func (C *CPU) txa() {
	switch C.Inst.addr {
	case implied:
		C.A = C.X
	default:
		log.Fatal("Bad addressing mode")
	}
	C.updateN(C.A)
	C.updateZ(C.A)

}

func (C *CPU) txs() {
	switch C.Inst.addr {
	case implied:
		C.SP = C.X
	default:
		log.Fatal("Bad addressing mode")
	}

}

func (C *CPU) tya() {
	switch C.Inst.addr {
	case implied:
		C.A = C.Y
	default:
		log.Fatal("Bad addressing mode")
	}
	C.updateN(C.A)
	C.updateZ(C.A)

}
