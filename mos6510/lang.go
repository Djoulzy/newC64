package mos6510

func (C *CPU) initLanguage() {
	mnemonic = map[byte]instruction{
		0xA9: {name: "LDA", addr: immediate, bytes: 2, cycles: 2, action: C.lda},
		0xA5: {name: "LDA", addr: zeropage, bytes: 2, cycles: 3, action: C.lda},
		0xB5: {name: "LDA", addr: zeropageX, bytes: 2, cycles: 4, action: C.lda},
		0xAD: {name: "LDA", addr: absolute, bytes: 3, cycles: 4, action: C.lda},
		0xBD: {name: "LDA", addr: absoluteX, bytes: 3, cycles: 4, action: C.lda},
		0xB9: {name: "LDA", addr: absoluteY, bytes: 3, cycles: 4, action: C.lda},
		0xA1: {name: "LDA", addr: indirectX, bytes: 2, cycles: 6, action: C.lda},
		0xB1: {name: "LDA", addr: indirectY, bytes: 2, cycles: 5, action: C.lda},
	}
}


func (C *CPU) lda() {

}