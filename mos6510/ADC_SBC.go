package mos6510

import (
	"log"

	"github.com/albenik/bcd"
)

// func bcd.ToUint8(i byte) byte {
// 	var bcd []byte
// 	for i > 0 {
// 		low := i % 10
// 		i /= 10
// 		hi := i % 10
// 		i /= 10
// 		var x []byte
// 		x = append(x, byte((hi&0xf)<<4)|byte(low&0xf))
// 		bcd = append(x, bcd[:]...)
// 	}
// 	if len(bcd) == 0 {
// 		return 0
// 	}
// 	return bcd[0]
// }

func (C *CPU) adc() {
	var val uint16
	var oper byte
	var crossed bool

	switch C.Inst.addr {
	case immediate:
		if C.getD() == 1 {
			A_bcd := uint16(bcd.ToUint8(C.A))
			oper_bcd := uint16(bcd.ToUint8(byte(C.oper)))
			val_bcd := A_bcd + oper_bcd + uint16(C.getC())
			vals := bcd.FromUint16(val_bcd)
			C.setC(vals[0] > 0)
			C.setV(val_bcd > 100)
			C.updateN(vals[1])
			C.updateZ(vals[1])
			C.A = vals[1]
			return
		} else {
			val = uint16(C.A) + C.oper + uint16(C.getC())
			C.setC(val > 0x00FF)
			C.updateV(C.A, byte(C.oper), byte(val))
			C.A = byte(val)
		}
	case zeropage:
		fallthrough
	case absolute:
		oper = C.ram.Read(C.oper)
		if C.getD() == 1 {
			A_bcd := uint16(bcd.ToUint8(C.A))
			oper_bcd := uint16(bcd.ToUint8(oper))
			val_bcd := A_bcd + oper_bcd + uint16(C.getC())
			vals := bcd.FromUint16(val_bcd)
			C.setC(vals[0] > 0)
			C.setV(val_bcd > 100)
			C.updateN(vals[1])
			C.updateZ(vals[1])
			C.A = vals[1]
			return
		} else {
			val = uint16(C.A) + uint16(oper) + uint16(C.getC())
		}
		C.setC(val > 0x00FF)
		C.updateV(C.A, oper, byte(val))
		C.A = byte(val)
	case zeropageX:
		oper = C.ram.Read(C.oper + uint16(C.X))
		if C.getD() == 1 {
			A_bcd := uint16(bcd.ToUint8(C.A))
			oper_bcd := uint16(bcd.ToUint8(oper))
			val_bcd := A_bcd + oper_bcd + uint16(C.getC())
			vals := bcd.FromUint16(val_bcd)
			C.setC(vals[0] > 0)
			C.setV(val_bcd > 100)
			C.updateN(vals[1])
			C.updateZ(vals[1])
			C.A = vals[1]
			return
		} else {
			val = uint16(C.A) + uint16(oper) + uint16(C.getC())
		}
		C.setC(val > 0x00FF)
		C.updateV(C.A, oper, byte(val))
		C.A = byte(val)
	case absoluteX:
		C.cross_oper = C.oper + uint16(C.X)
		if C.oper&0xFF00 == C.cross_oper&0xFF00 {
			oper = C.ram.Read(C.cross_oper)
			if C.getD() == 1 {
				A_bcd := uint16(bcd.ToUint8(C.A))
				oper_bcd := uint16(bcd.ToUint8(oper))
				val_bcd := A_bcd + oper_bcd + uint16(C.getC())
				vals := bcd.FromUint16(val_bcd)
				C.setC(vals[0] > 0)
				C.setV(val_bcd > 100)
				C.updateN(vals[1])
				C.updateZ(vals[1])
				C.A = vals[1]
				return
			} else {
				val = uint16(C.A) + uint16(oper) + uint16(C.getC())
			}
			C.setC(val > 0x00FF)
			C.updateV(C.A, oper, byte(val))
			C.A = byte(val)
		} else {
			C.Inst.addr = CrossPage
			C.State = Compute
			C.Inst.Cycles++
			return
		}
	case absoluteY:
		C.cross_oper = C.oper + uint16(C.Y)
		if C.oper&0xFF00 == C.cross_oper&0xFF00 {
			oper = C.ram.Read(C.cross_oper)
			if C.getD() == 1 {
				A_bcd := uint16(bcd.ToUint8(C.A))
				oper_bcd := uint16(bcd.ToUint8(oper))
				val_bcd := A_bcd + oper_bcd + uint16(C.getC())
				vals := bcd.FromUint16(val_bcd)
				C.setC(vals[0] > 0)
				C.setV(val_bcd > 100)
				C.updateN(vals[1])
				C.updateZ(vals[1])
				C.A = vals[1]
				return
			} else {
				val = uint16(C.A) + uint16(oper) + uint16(C.getC())
			}
			C.setC(val > 0x00FF)
			C.updateV(C.A, oper, byte(val))
			C.A = byte(val)
		} else {
			C.Inst.addr = CrossPage
			C.State = Compute
			C.Inst.Cycles++
			return
		}
	case indirectX:
		oper = C.ReadIndirectX(C.oper)
		if C.getD() == 1 {
			A_bcd := uint16(bcd.ToUint8(C.A))
			oper_bcd := uint16(bcd.ToUint8(oper))
			val_bcd := A_bcd + oper_bcd + uint16(C.getC())
			vals := bcd.FromUint16(val_bcd)
			C.setC(vals[0] > 0)
			C.setV(val_bcd > 100)
			C.updateN(vals[1])
			C.updateZ(vals[1])
			C.A = vals[1]
			return
		} else {
			val = uint16(C.A) + uint16(oper) + uint16(C.getC())
		}
		C.setC(val > 0x00FF)
		C.updateV(C.A, oper, byte(val))
		C.A = byte(val)
	case indirectY:
		C.cross_oper = C.GetIndirectYAddr(C.oper, &crossed)
		if crossed {
			oper = C.ram.Read(C.cross_oper)
			if C.getD() == 1 {
				A_bcd := uint16(bcd.ToUint8(C.A))
				oper_bcd := uint16(bcd.ToUint8(oper))
				val_bcd := A_bcd + oper_bcd + uint16(C.getC())
				vals := bcd.FromUint16(val_bcd)
				C.setC(vals[0] > 0)
				C.setV(val_bcd > 100)
				C.updateN(vals[1])
				C.updateZ(vals[1])
				C.A = vals[1]
				return
			} else {
				val = uint16(C.A) + uint16(oper) + uint16(C.getC())
			}
			C.setC(val > 0x00FF)
			C.updateV(C.A, oper, byte(val))
			C.A = byte(val)
		} else {
			C.Inst.addr = CrossPage
			C.State = Compute
			C.Inst.Cycles++
			return
		}
	case CrossPage:
		oper = C.ram.Read(C.cross_oper)
		if C.getD() == 1 {
			A_bcd := uint16(bcd.ToUint8(C.A))
			oper_bcd := uint16(bcd.ToUint8(oper))
			val_bcd := A_bcd + oper_bcd + uint16(C.getC())
			vals := bcd.FromUint16(val_bcd)
			C.setC(vals[0] > 0)
			C.setV(val_bcd > 100)
			C.updateN(vals[1])
			C.updateZ(vals[1])
			C.A = vals[1]
			return
		} else {
			val = uint16(C.A) + uint16(oper) + uint16(C.getC())
		}
		C.setC(val > 0x00FF)
		C.updateV(C.A, oper, byte(val))
		C.A = byte(val)
	default:
		log.Fatal("Bad addressing mode")
	}
	C.updateN(byte(val))
	C.updateZ(byte(val))
}

func (C *CPU) sbc() {
	var val int
	var oper byte
	var crossed bool

	switch C.Inst.addr {
	case immediate:
		if C.getD() == 1 {
			A_bcd := uint16(bcd.ToUint8(C.A))
			oper_bcd := uint16(bcd.ToUint8(byte(C.oper)))
			val_bcd := A_bcd - oper_bcd
			if C.getC() == 0 {
				val_bcd -= 1
			}
			vals := bcd.FromUint16(val_bcd)
			C.setC(vals[1] > 0)
			C.setV(val_bcd > 100)
			C.updateN(vals[1])
			C.updateZ(vals[1])
			C.A = vals[1]
			return
		} else {
			val = int(C.A) - int(C.oper)
		}
		if C.getC() == 0 {
			val -= 1
		}
		C.updateV(C.A, ^byte(C.oper), byte(val))
		C.A = byte(val)
	case zeropage:
		fallthrough
	case absolute:
		oper = C.ram.Read(C.oper)
		if C.getD() == 1 {
			A_bcd := uint16(bcd.ToUint8(C.A))
			oper_bcd := uint16(bcd.ToUint8(byte(oper)))
			val_bcd := A_bcd - oper_bcd
			if C.getC() == 0 {
				val_bcd -= 1
			}
			vals := bcd.FromUint16(val_bcd)
			C.setC(vals[1] > 0)
			C.setV(val_bcd > 100)
			C.updateN(vals[1])
			C.updateZ(vals[1])
			C.A = vals[1]
			return
		} else {
			val = int(C.A) - int(oper)
		}
		if C.getC() == 0 {
			val -= 1
		}
		C.updateV(C.A, ^oper, byte(val))
		C.A = byte(val)
	case zeropageX:
		oper = C.ram.Read(C.oper + uint16(C.X))
		if C.getD() == 1 {
			A_bcd := uint16(bcd.ToUint8(C.A))
			oper_bcd := uint16(bcd.ToUint8(byte(oper)))
			val_bcd := A_bcd - oper_bcd
			if C.getC() == 0 {
				val_bcd -= 1
			}
			vals := bcd.FromUint16(val_bcd)
			C.setC(vals[1] > 0)
			C.setV(val_bcd > 100)
			C.updateN(vals[1])
			C.updateZ(vals[1])
			C.A = vals[1]
			return
		} else {
			val = int(C.A) - int(oper)
		}
		if C.getC() == 0 {
			val -= 1
		}
		C.updateV(C.A, ^oper, byte(val))
		C.A = byte(val)
	case absoluteX:
		C.cross_oper = C.oper + uint16(C.X)
		if C.oper&0xFF00 == C.cross_oper&0xFF00 {
			oper = C.ram.Read(C.cross_oper)
			if C.getD() == 1 {
				A_bcd := uint16(bcd.ToUint8(C.A))
				oper_bcd := uint16(bcd.ToUint8(byte(oper)))
				val_bcd := A_bcd - oper_bcd
				if C.getC() == 0 {
					val_bcd -= 1
				}
				vals := bcd.FromUint16(val_bcd)
				C.setC(vals[1] > 0)
				C.setV(val_bcd > 100)
				C.updateN(vals[1])
				C.updateZ(vals[1])
				C.A = vals[1]
				return
			} else {
				val = int(C.A) - int(oper)
			}
			if C.getC() == 0 {
				val -= 1
			}
			C.updateV(C.A, ^oper, byte(val))
			C.A = byte(val)
		} else {
			C.Inst.addr = CrossPage
			C.State = Compute
			C.Inst.Cycles++
			return
		}
	case absoluteY:
		C.cross_oper = C.oper + uint16(C.Y)
		if C.oper&0xFF00 == C.cross_oper&0xFF00 {
			oper = C.ram.Read(C.cross_oper)
			if C.getD() == 1 {
				A_bcd := uint16(bcd.ToUint8(C.A))
				oper_bcd := uint16(bcd.ToUint8(byte(oper)))
				val_bcd := A_bcd - oper_bcd
				if C.getC() == 0 {
					val_bcd -= 1
				}
				vals := bcd.FromUint16(val_bcd)
				C.setC(vals[1] > 0)
				C.setV(val_bcd > 100)
				C.updateN(vals[1])
				C.updateZ(vals[1])
				C.A = vals[1]
				return
			} else {
				val = int(C.A) - int(oper)
			}
			if C.getC() == 0 {
				val -= 1
			}
			C.updateV(C.A, ^oper, byte(val))
			C.A = byte(val)
		} else {
			C.Inst.addr = CrossPage
			C.State = Compute
			C.Inst.Cycles++
			return
		}
	case indirectX:
		oper = C.ReadIndirectX(C.oper)
		if C.getD() == 1 {
			A_bcd := uint16(bcd.ToUint8(C.A))
			oper_bcd := uint16(bcd.ToUint8(byte(oper)))
			val_bcd := A_bcd - oper_bcd
			if C.getC() == 0 {
				val_bcd -= 1
			}
			vals := bcd.FromUint16(val_bcd)
			C.setC(vals[1] > 0)
			C.setV(val_bcd > 100)
			C.updateN(vals[1])
			C.updateZ(vals[1])
			C.A = vals[1]
			return
		} else {
			val = int(C.A) - int(oper)
		}
		if C.getC() == 0 {
			val -= 1
		}
		C.updateV(C.A, ^oper, byte(val))
		C.A = byte(val)
	case indirectY:
		C.cross_oper = C.GetIndirectYAddr(C.oper, &crossed)
		if crossed {
			oper = C.ram.Read(C.cross_oper)
			if C.getD() == 1 {
				A_bcd := uint16(bcd.ToUint8(C.A))
				oper_bcd := uint16(bcd.ToUint8(byte(oper)))
				val_bcd := A_bcd - oper_bcd
				if C.getC() == 0 {
					val_bcd -= 1
				}
				vals := bcd.FromUint16(val_bcd)
				C.setC(vals[1] > 0)
				C.setV(val_bcd > 100)
				C.updateN(vals[1])
				C.updateZ(vals[1])
				C.A = vals[1]
				return
			} else {
				val = int(C.A) - int(oper)
			}
			if C.getC() == 0 {
				val -= 1
			}
			C.updateV(C.A, ^oper, byte(val))
			C.A = byte(val)
		} else {
			C.Inst.addr = CrossPage
			C.State = Compute
			C.Inst.Cycles++
			return
		}
	case CrossPage:
		oper = C.ram.Read(C.cross_oper)
		if C.getD() == 1 {
			A_bcd := uint16(bcd.ToUint8(C.A))
			oper_bcd := uint16(bcd.ToUint8(byte(oper)))
			val_bcd := A_bcd - oper_bcd
			if C.getC() == 0 {
				val_bcd -= 1
			}
			vals := bcd.FromUint16(val_bcd)
			C.setC(vals[1] > 0)
			C.setV(val_bcd > 100)
			C.updateN(vals[1])
			C.updateZ(vals[1])
			C.A = vals[1]
			return
		} else {
			val = int(C.A) - int(oper)
		}
		if C.getC() == 0 {
			val -= 1
		}
		C.updateV(C.A, ^oper, byte(val))
		C.A = byte(val)
	default:
		log.Fatal("Bad addressing mode")
	}
	C.setC(val >= 0x00)
	C.setN(val&0b10000000 == 0b10000000)
	C.updateZ(byte(val))
}
