package mos6510

func (C *CPU) initLanguage() {
	mnemonic = map[byte]instruction{

		0x69: {name: "ADC", bytes: 2, Cycles: 2, action: C.adc, addr: immediate},
		0x65: {name: "ADC", bytes: 2, Cycles: 3, action: C.adc, addr: zeropage},
		0x75: {name: "ADC", bytes: 2, Cycles: 4, action: C.adc, addr: zeropageX},
		0x6D: {name: "ADC", bytes: 3, Cycles: 4, action: C.adc, addr: absolute},
		0x7D: {name: "ADC", bytes: 3, Cycles: 4, action: C.adc, addr: absoluteX},
		0x79: {name: "ADC", bytes: 3, Cycles: 4, action: C.adc, addr: absoluteY},
		0x61: {name: "ADC", bytes: 2, Cycles: 6, action: C.adc, addr: indirectX},
		0x71: {name: "ADC", bytes: 2, Cycles: 5, action: C.adc, addr: indirectY},

		0x0B: {name: "ANC", bytes: 2, Cycles: 2, action: C.anc, addr: immediate},

		0x29: {name: "AND", bytes: 2, Cycles: 2, action: C.and, addr: immediate},
		0x25: {name: "AND", bytes: 2, Cycles: 3, action: C.and, addr: zeropage},
		0x35: {name: "AND", bytes: 2, Cycles: 4, action: C.and, addr: zeropageX},
		0x2D: {name: "AND", bytes: 3, Cycles: 4, action: C.and, addr: absolute},
		0x3D: {name: "AND", bytes: 3, Cycles: 4, action: C.and, addr: absoluteX},
		0x39: {name: "AND", bytes: 3, Cycles: 4, action: C.and, addr: absoluteY},
		0x21: {name: "AND", bytes: 2, Cycles: 6, action: C.and, addr: indirectX},
		0x31: {name: "AND", bytes: 2, Cycles: 5, action: C.and, addr: indirectY},

		0x4B: {name: "ALR", bytes: 2, Cycles: 2, action: C.alr, addr: immediate},

		0x0A: {name: "ASL", bytes: 1, Cycles: 2, action: C.asl, addr: implied},
		0x06: {name: "ASL", bytes: 2, Cycles: 5, action: C.asl, addr: zeropage},
		0x16: {name: "ASL", bytes: 2, Cycles: 6, action: C.asl, addr: zeropageX},
		0x0E: {name: "ASL", bytes: 3, Cycles: 6, action: C.asl, addr: absolute},
		0x1E: {name: "ASL", bytes: 3, Cycles: 7, action: C.asl, addr: absoluteX},

		0x90: {name: "BCC", bytes: 2, Cycles: 2, action: C.bcc, addr: relative},

		0xB0: {name: "BCS", bytes: 2, Cycles: 2, action: C.bcs, addr: relative},

		0xF0: {name: "BEQ", bytes: 2, Cycles: 2, action: C.beq, addr: relative},

		0x24: {name: "BIT", bytes: 2, Cycles: 3, action: C.bit, addr: zeropage},
		0x2C: {name: "BIT", bytes: 3, Cycles: 4, action: C.bit, addr: absolute},

		0x30: {name: "BMI", bytes: 2, Cycles: 2, action: C.bmi, addr: relative},

		0xD0: {name: "BNE", bytes: 2, Cycles: 2, action: C.bne, addr: relative},

		0x10: {name: "BPL", bytes: 2, Cycles: 2, action: C.bpl, addr: relative},

		0x00: {name: "BRK", bytes: 1, Cycles: 7, action: C.brk, addr: implied},

		0x50: {name: "BVC", bytes: 2, Cycles: 2, action: C.bvc, addr: relative},

		0x70: {name: "BVS", bytes: 2, Cycles: 2, action: C.bvs, addr: relative},

		0x18: {name: "CLC", bytes: 1, Cycles: 2, action: C.clc, addr: implied},

		0xD8: {name: "CLD", bytes: 1, Cycles: 2, action: C.cld, addr: implied},

		0x58: {name: "CLI", bytes: 1, Cycles: 2, action: C.cli, addr: implied},

		0xB8: {name: "CLV", bytes: 1, Cycles: 2, action: C.clv, addr: implied},

		0xC9: {name: "CMP", bytes: 2, Cycles: 2, action: C.cmp, addr: immediate},
		0xC5: {name: "CMP", bytes: 2, Cycles: 3, action: C.cmp, addr: zeropage},
		0xD5: {name: "CMP", bytes: 2, Cycles: 4, action: C.cmp, addr: zeropageX},
		0xCD: {name: "CMP", bytes: 3, Cycles: 4, action: C.cmp, addr: absolute},
		0xDD: {name: "CMP", bytes: 3, Cycles: 4, action: C.cmp, addr: absoluteX},
		0xD9: {name: "CMP", bytes: 3, Cycles: 4, action: C.cmp, addr: absoluteY},
		0xC1: {name: "CMP", bytes: 2, Cycles: 6, action: C.cmp, addr: indirectX},
		0xD1: {name: "CMP", bytes: 2, Cycles: 5, action: C.cmp, addr: indirectY},

		0xE0: {name: "CPX", bytes: 2, Cycles: 2, action: C.cpx, addr: immediate},
		0xE4: {name: "CPX", bytes: 2, Cycles: 3, action: C.cpx, addr: zeropage},
		0xEC: {name: "CPX", bytes: 3, Cycles: 4, action: C.cpx, addr: absolute},

		0xC0: {name: "CPY", bytes: 2, Cycles: 2, action: C.cpy, addr: immediate},
		0xC4: {name: "CPY", bytes: 2, Cycles: 3, action: C.cpy, addr: zeropage},
		0xCC: {name: "CPY", bytes: 3, Cycles: 4, action: C.cpy, addr: absolute},

		0xC7: {name: "DCP", bytes: 2, Cycles: 5, action: C.dcp, addr: zeropage},
		0xD7: {name: "DCP", bytes: 2, Cycles: 6, action: C.dcp, addr: zeropageX},
		0xC3: {name: "DCP", bytes: 2, Cycles: 8, action: C.dcp, addr: indirectX},
		0xD3: {name: "DCP", bytes: 2, Cycles: 8, action: C.dcp, addr: indirectY},
		0xCF: {name: "DCP", bytes: 3, Cycles: 6, action: C.dcp, addr: absolute},
		0xDF: {name: "DCP", bytes: 3, Cycles: 7, action: C.dcp, addr: absoluteX},
		0xDB: {name: "DCP", bytes: 3, Cycles: 7, action: C.dcp, addr: absoluteY},

		0xC6: {name: "DEC", bytes: 2, Cycles: 5, action: C.dec, addr: zeropage},
		0xD6: {name: "DEC", bytes: 2, Cycles: 6, action: C.dec, addr: zeropageX},
		0xCE: {name: "DEC", bytes: 3, Cycles: 6, action: C.dec, addr: absolute},
		0xDE: {name: "DEC", bytes: 3, Cycles: 7, action: C.dec, addr: absoluteX},

		0xCA: {name: "DEX", bytes: 1, Cycles: 2, action: C.dex, addr: implied},

		0x88: {name: "DEY", bytes: 1, Cycles: 2, action: C.dey, addr: implied},

		0x49: {name: "EOR", bytes: 2, Cycles: 2, action: C.eor, addr: immediate},
		0x45: {name: "EOR", bytes: 2, Cycles: 3, action: C.eor, addr: zeropage},
		0x55: {name: "EOR", bytes: 2, Cycles: 4, action: C.eor, addr: zeropageX},
		0x4D: {name: "EOR", bytes: 3, Cycles: 4, action: C.eor, addr: absolute},
		0x5D: {name: "EOR", bytes: 3, Cycles: 4, action: C.eor, addr: absoluteX},
		0x59: {name: "EOR", bytes: 3, Cycles: 4, action: C.eor, addr: absoluteY},
		0x41: {name: "EOR", bytes: 2, Cycles: 6, action: C.eor, addr: indirectX},
		0x51: {name: "EOR", bytes: 2, Cycles: 5, action: C.eor, addr: indirectY},

		0xE6: {name: "INC", bytes: 2, Cycles: 5, action: C.inc, addr: zeropage},
		0xF6: {name: "INC", bytes: 2, Cycles: 6, action: C.inc, addr: zeropageX},
		0xEE: {name: "INC", bytes: 3, Cycles: 6, action: C.inc, addr: absolute},
		0xFE: {name: "INC", bytes: 3, Cycles: 7, action: C.inc, addr: absoluteX},

		0xE8: {name: "INX", bytes: 1, Cycles: 2, action: C.inx, addr: implied},

		0xC8: {name: "INY", bytes: 1, Cycles: 2, action: C.iny, addr: implied},

		0xE7: {name: "ISC", bytes: 2, Cycles: 5, action: C.isc, addr: zeropage},
		0xF7: {name: "ISC", bytes: 2, Cycles: 6, action: C.isc, addr: zeropageX},
		0xE3: {name: "ISC", bytes: 2, Cycles: 8, action: C.isc, addr: indirectX},
		0xF3: {name: "ISC", bytes: 2, Cycles: 8, action: C.isc, addr: indirectY},
		0xEF: {name: "ISC", bytes: 3, Cycles: 6, action: C.isc, addr: absolute},
		0xFF: {name: "ISC", bytes: 3, Cycles: 7, action: C.isc, addr: absoluteX},
		0xFB: {name: "ISC", bytes: 3, Cycles: 7, action: C.isc, addr: absoluteY},

		0x4C: {name: "JMP", bytes: 3, Cycles: 3, action: C.jmp, addr: absolute},
		0x6C: {name: "JMP", bytes: 3, Cycles: 5, action: C.jmp, addr: indirect},

		0x20: {name: "JSR", bytes: 3, Cycles: 6, action: C.jsr, addr: absolute},

		0x02: {name: "KIL", bytes: 1, Cycles: 1, action: func() { C.State = Idle }, addr: implied},
		0x12: {name: "KIL", bytes: 1, Cycles: 1, action: func() { C.State = Idle }, addr: implied},
		0x22: {name: "KIL", bytes: 1, Cycles: 1, action: func() { C.State = Idle }, addr: implied},
		0x32: {name: "KIL", bytes: 1, Cycles: 1, action: func() { C.State = Idle }, addr: implied},
		0x42: {name: "KIL", bytes: 1, Cycles: 1, action: func() { C.State = Idle }, addr: implied},
		0x52: {name: "KIL", bytes: 1, Cycles: 1, action: func() { C.State = Idle }, addr: implied},
		0x62: {name: "KIL", bytes: 1, Cycles: 1, action: func() { C.State = Idle }, addr: implied},
		0x72: {name: "KIL", bytes: 1, Cycles: 1, action: func() { C.State = Idle }, addr: implied},
		0x92: {name: "KIL", bytes: 1, Cycles: 1, action: func() { C.State = Idle }, addr: implied},
		0xB2: {name: "KIL", bytes: 1, Cycles: 1, action: func() { C.State = Idle }, addr: implied},
		0xD2: {name: "KIL", bytes: 1, Cycles: 1, action: func() { C.State = Idle }, addr: implied},
		0xF2: {name: "KIL", bytes: 1, Cycles: 1, action: func() { C.State = Idle }, addr: implied},

		0xA9: {name: "LDA", bytes: 2, Cycles: 2, action: C.lda, addr: immediate},
		0xA5: {name: "LDA", bytes: 2, Cycles: 3, action: C.lda, addr: zeropage},
		0xB5: {name: "LDA", bytes: 2, Cycles: 4, action: C.lda, addr: zeropageX},
		0xAD: {name: "LDA", bytes: 3, Cycles: 4, action: C.lda, addr: absolute},
		0xBD: {name: "LDA", bytes: 3, Cycles: 4, action: C.lda, addr: absoluteX},
		0xB9: {name: "LDA", bytes: 3, Cycles: 4, action: C.lda, addr: absoluteY},
		0xA1: {name: "LDA", bytes: 2, Cycles: 6, action: C.lda, addr: indirectX},
		0xB1: {name: "LDA", bytes: 2, Cycles: 5, action: C.lda, addr: indirectY},

		0xA2: {name: "LDX", bytes: 2, Cycles: 2, action: C.ldx, addr: immediate},
		0xA6: {name: "LDX", bytes: 2, Cycles: 3, action: C.ldx, addr: zeropage},
		0xB6: {name: "LDX", bytes: 2, Cycles: 4, action: C.ldx, addr: zeropageY},
		0xAE: {name: "LDX", bytes: 3, Cycles: 4, action: C.ldx, addr: absolute},
		0xBE: {name: "LDX", bytes: 3, Cycles: 4, action: C.ldx, addr: absoluteY},

		0xA0: {name: "LDY", bytes: 2, Cycles: 2, action: C.ldy, addr: immediate},
		0xA4: {name: "LDY", bytes: 2, Cycles: 3, action: C.ldy, addr: zeropage},
		0xB4: {name: "LDY", bytes: 2, Cycles: 4, action: C.ldy, addr: zeropageX},
		0xAC: {name: "LDY", bytes: 3, Cycles: 4, action: C.ldy, addr: absolute},
		0xBC: {name: "LDY", bytes: 3, Cycles: 4, action: C.ldy, addr: absoluteX},

		0x4A: {name: "LSR", bytes: 1, Cycles: 2, action: C.lsr, addr: implied},
		0x46: {name: "LSR", bytes: 2, Cycles: 5, action: C.lsr, addr: zeropage},
		0x56: {name: "LSR", bytes: 2, Cycles: 6, action: C.lsr, addr: zeropageX},
		0x4E: {name: "LSR", bytes: 3, Cycles: 6, action: C.lsr, addr: absolute},
		0x5E: {name: "LSR", bytes: 3, Cycles: 7, action: C.lsr, addr: absoluteX},

		0xEA: {name: "NOP", bytes: 1, Cycles: 2, action: C.nop, addr: implied},
		0x1A: {name: "NOP", bytes: 1, Cycles: 2, action: C.nop, addr: implied},
		0x3A: {name: "NOP", bytes: 1, Cycles: 2, action: C.nop, addr: implied},
		0x5A: {name: "NOP", bytes: 1, Cycles: 2, action: C.nop, addr: implied},
		0x7A: {name: "NOP", bytes: 1, Cycles: 2, action: C.nop, addr: implied},
		0xDA: {name: "NOP", bytes: 1, Cycles: 2, action: C.nop, addr: implied},
		0xFA: {name: "NOP", bytes: 1, Cycles: 2, action: C.nop, addr: implied},
		0x80: {name: "NOP", bytes: 2, Cycles: 2, action: C.nop, addr: immediate},
		0x82: {name: "NOP", bytes: 2, Cycles: 2, action: C.nop, addr: immediate},
		0xC2: {name: "NOP", bytes: 2, Cycles: 2, action: C.nop, addr: immediate},
		0xE2: {name: "NOP", bytes: 2, Cycles: 2, action: C.nop, addr: immediate},
		0x89: {name: "NOP", bytes: 2, Cycles: 2, action: C.nop, addr: immediate},
		0x04: {name: "NOP", bytes: 2, Cycles: 3, action: C.nop, addr: zeropage},
		0x44: {name: "NOP", bytes: 2, Cycles: 3, action: C.nop, addr: zeropage},
		0x64: {name: "NOP", bytes: 2, Cycles: 3, action: C.nop, addr: zeropage},
		0x14: {name: "NOP", bytes: 2, Cycles: 4, action: C.nop, addr: zeropageX},
		0x34: {name: "NOP", bytes: 2, Cycles: 4, action: C.nop, addr: zeropageX},
		0x54: {name: "NOP", bytes: 2, Cycles: 4, action: C.nop, addr: zeropageX},
		0x74: {name: "NOP", bytes: 2, Cycles: 4, action: C.nop, addr: zeropageX},
		0xD4: {name: "NOP", bytes: 2, Cycles: 4, action: C.nop, addr: zeropageX},
		0xF4: {name: "NOP", bytes: 2, Cycles: 4, action: C.nop, addr: zeropageX},
		0x0C: {name: "NOP", bytes: 3, Cycles: 4, action: C.nop, addr: absolute},
		0x1C: {name: "NOP", bytes: 3, Cycles: 4, action: C.nop, addr: absoluteX},
		0x3C: {name: "NOP", bytes: 3, Cycles: 4, action: C.nop, addr: absoluteX},
		0x5C: {name: "NOP", bytes: 3, Cycles: 4, action: C.nop, addr: absoluteX},
		0x7C: {name: "NOP", bytes: 3, Cycles: 4, action: C.nop, addr: absoluteX},
		0xDC: {name: "NOP", bytes: 3, Cycles: 4, action: C.nop, addr: absoluteX},
		0xFC: {name: "NOP", bytes: 3, Cycles: 4, action: C.nop, addr: absoluteX},

		0x09: {name: "ORA", bytes: 2, Cycles: 2, action: C.ora, addr: immediate},
		0x05: {name: "ORA", bytes: 2, Cycles: 3, action: C.ora, addr: zeropage},
		0x15: {name: "ORA", bytes: 2, Cycles: 4, action: C.ora, addr: zeropageX},
		0x0D: {name: "ORA", bytes: 3, Cycles: 4, action: C.ora, addr: absolute},
		0x1D: {name: "ORA", bytes: 3, Cycles: 4, action: C.ora, addr: absoluteX},
		0x19: {name: "ORA", bytes: 3, Cycles: 4, action: C.ora, addr: absoluteY},
		0x01: {name: "ORA", bytes: 2, Cycles: 6, action: C.ora, addr: indirectX},
		0x11: {name: "ORA", bytes: 2, Cycles: 5, action: C.ora, addr: indirectY},

		0x48: {name: "PHA", bytes: 1, Cycles: 3, action: C.pha, addr: implied},

		0x08: {name: "PHP", bytes: 1, Cycles: 3, action: C.php, addr: implied},

		0x68: {name: "PLA", bytes: 1, Cycles: 4, action: C.pla, addr: implied},

		0x28: {name: "PLP", bytes: 1, Cycles: 4, action: C.plp, addr: implied},

		0x27: {name: "RLA", bytes: 2, Cycles: 5, action: C.rla, addr: zeropage},
		0x37: {name: "RLA", bytes: 2, Cycles: 6, action: C.rla, addr: zeropageX},
		0x23: {name: "RLA", bytes: 2, Cycles: 8, action: C.rla, addr: indirectX},
		0x33: {name: "RLA", bytes: 2, Cycles: 8, action: C.rla, addr: indirectY},
		0x2F: {name: "RLA", bytes: 3, Cycles: 6, action: C.rla, addr: absolute},
		0x3F: {name: "RLA", bytes: 3, Cycles: 7, action: C.rla, addr: absoluteX},
		0x3B: {name: "RLA", bytes: 3, Cycles: 7, action: C.rla, addr: absoluteY},

		0x2A: {name: "ROL", bytes: 1, Cycles: 2, action: C.rol, addr: implied},
		0x26: {name: "ROL", bytes: 2, Cycles: 5, action: C.rol, addr: zeropage},
		0x36: {name: "ROL", bytes: 2, Cycles: 6, action: C.rol, addr: zeropageX},
		0x2E: {name: "ROL", bytes: 3, Cycles: 6, action: C.rol, addr: absolute},
		0x3E: {name: "ROL", bytes: 3, Cycles: 7, action: C.rol, addr: absoluteX},

		0x6A: {name: "ROR", bytes: 1, Cycles: 2, action: C.ror, addr: implied},
		0x66: {name: "ROR", bytes: 2, Cycles: 5, action: C.ror, addr: zeropage},
		0x76: {name: "ROR", bytes: 2, Cycles: 6, action: C.ror, addr: zeropageX},
		0x6E: {name: "ROR", bytes: 3, Cycles: 6, action: C.ror, addr: absolute},
		0x7E: {name: "ROR", bytes: 3, Cycles: 7, action: C.ror, addr: absoluteX},

		0x40: {name: "RTI", bytes: 1, Cycles: 6, action: C.rti, addr: implied},

		0x60: {name: "RTS", bytes: 1, Cycles: 6, action: C.rts, addr: implied},

		0x87: {name: "SAX", bytes: 2, Cycles: 3, action: C.sax, addr: zeropage},
		0x97: {name: "SAX", bytes: 2, Cycles: 4, action: C.sax, addr: zeropageY},
		0x83: {name: "SAX", bytes: 2, Cycles: 6, action: C.sax, addr: zeropageX},
		0x8F: {name: "SAX", bytes: 3, Cycles: 4, action: C.sax, addr: absolute},

		0xE9: {name: "SBC", bytes: 2, Cycles: 2, action: C.sbc, addr: immediate},
		0xE5: {name: "SBC", bytes: 2, Cycles: 3, action: C.sbc, addr: zeropage},
		0xF5: {name: "SBC", bytes: 2, Cycles: 4, action: C.sbc, addr: zeropageX},
		0xED: {name: "SBC", bytes: 3, Cycles: 4, action: C.sbc, addr: absolute},
		0xFD: {name: "SBC", bytes: 3, Cycles: 4, action: C.sbc, addr: absoluteX},
		0xF9: {name: "SBC", bytes: 3, Cycles: 4, action: C.sbc, addr: absoluteY},
		0xE1: {name: "SBC", bytes: 2, Cycles: 6, action: C.sbc, addr: indirectX},
		0xF1: {name: "SBC", bytes: 2, Cycles: 5, action: C.sbc, addr: indirectY},

		0xCB: {name: "SBX", bytes: 2, Cycles: 2, action: C.sbx, addr: immediate},

		0x07: {name: "SLO", bytes: 2, Cycles: 5, action: C.slo, addr: zeropage},
		0x17: {name: "SLO", bytes: 2, Cycles: 6, action: C.slo, addr: zeropageX},
		0x03: {name: "SLO", bytes: 2, Cycles: 8, action: C.slo, addr: indirectX},
		0x13: {name: "SLO", bytes: 2, Cycles: 8, action: C.slo, addr: indirectY},
		0x0F: {name: "SLO", bytes: 3, Cycles: 6, action: C.slo, addr: absolute},
		0x1F: {name: "SLO", bytes: 3, Cycles: 7, action: C.slo, addr: absoluteX},
		0x1B: {name: "SLO", bytes: 3, Cycles: 7, action: C.slo, addr: absoluteY},

		0x38: {name: "SEC", bytes: 1, Cycles: 2, action: C.sec, addr: implied},

		0xF8: {name: "SED", bytes: 1, Cycles: 2, action: C.sed, addr: implied},

		0x78: {name: "SEI", bytes: 1, Cycles: 2, action: C.sei, addr: implied},

		0x47: {name: "SRE", bytes: 2, Cycles: 5, action: C.sre, addr: zeropage},
		0x57: {name: "SRE", bytes: 2, Cycles: 6, action: C.sre, addr: zeropageX},
		0x43: {name: "SRE", bytes: 2, Cycles: 8, action: C.sre, addr: indirectX},
		0x53: {name: "SRE", bytes: 2, Cycles: 8, action: C.sre, addr: indirectY},
		0x4F: {name: "SRE", bytes: 3, Cycles: 6, action: C.sre, addr: absolute},
		0x5F: {name: "SRE", bytes: 3, Cycles: 7, action: C.sre, addr: absoluteX},
		0x5B: {name: "SRE", bytes: 3, Cycles: 7, action: C.sre, addr: absoluteY},

		0x85: {name: "STA", bytes: 2, Cycles: 3, action: C.sta, addr: zeropage},
		0x95: {name: "STA", bytes: 2, Cycles: 4, action: C.sta, addr: zeropageX},
		0x8D: {name: "STA", bytes: 3, Cycles: 4, action: C.sta, addr: absolute},
		0x9D: {name: "STA", bytes: 3, Cycles: 5, action: C.sta, addr: absoluteX},
		0x99: {name: "STA", bytes: 3, Cycles: 5, action: C.sta, addr: absoluteY},
		0x81: {name: "STA", bytes: 2, Cycles: 6, action: C.sta, addr: indirectX},
		0x91: {name: "STA", bytes: 2, Cycles: 6, action: C.sta, addr: indirectY},

		0x86: {name: "STX", bytes: 2, Cycles: 3, action: C.stx, addr: zeropage},
		0x96: {name: "STX", bytes: 2, Cycles: 4, action: C.stx, addr: zeropageY},
		0x8E: {name: "STX", bytes: 3, Cycles: 4, action: C.stx, addr: absolute},

		0x84: {name: "STY", bytes: 2, Cycles: 3, action: C.sty, addr: zeropage},
		0x94: {name: "STY", bytes: 2, Cycles: 4, action: C.sty, addr: zeropageX},
		0x8C: {name: "STY", bytes: 3, Cycles: 4, action: C.sty, addr: absolute},

		0xAA: {name: "TAX", bytes: 1, Cycles: 2, action: C.tax, addr: implied},

		0xA8: {name: "TAY", bytes: 1, Cycles: 2, action: C.tay, addr: implied},

		0xBA: {name: "TSX", bytes: 1, Cycles: 2, action: C.tsx, addr: implied},

		0x8A: {name: "TXA", bytes: 1, Cycles: 2, action: C.txa, addr: implied},

		0x9A: {name: "TXS", bytes: 1, Cycles: 2, action: C.txs, addr: implied},

		0x98: {name: "TYA", bytes: 1, Cycles: 2, action: C.tya, addr: implied},
	}
}
