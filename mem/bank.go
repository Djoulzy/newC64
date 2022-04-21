package mem

import (
	"fmt"
	"newC64/clog"
	"newC64/trace"
)

type BANK struct {
	Selector *byte
	Layouts  []CONFIG
}

func InitBanks(nbMemLayout int, sel *byte) BANK {
	B := BANK{}
	B.Layouts = make([]CONFIG, nbMemLayout)
	B.Selector = sel
	return B
}

func (B *BANK) Read(addr uint16) byte {
	// clog.Test("MEM", "Read", "Addr: %04X, Page: %d, Selector: %d", addr, int(addr>>PAGE_DIVIDER), *B.Selector&0x1F)
	bank := B.Layouts[*B.Selector&0x1F]
	layerNum := bank.LayerByPages[int(addr>>PAGE_DIVIDER)]
	// return C.Layers[layerNum][addr-C.Start[layerNum]]
	// clog.Test("MEM", "Read", "Addr: %04X, Page: %d, Layer: %d", addr, int(addr>>PAGE_DIVIDER), layerNum)
	return bank.Accessors[layerNum].MRead(bank.Layers[layerNum], addr-bank.Start[layerNum])
}

func (B *BANK) Write(addr uint16, value byte) {
	bank := B.Layouts[*B.Selector&0x1F]
	layerNum := bank.LayerByPages[int(addr>>PAGE_DIVIDER)]
	if bank.ReadOnly[layerNum] {
		layerNum = 0
	}
	// clog.Test("MEM", "Write", "Addr: %04X, Page: %d, Layer: %d", addr, int(addr>>PAGE_DIVIDER), layerNum)
	bank.Accessors[layerNum].MWrite(bank.Layers[layerNum], addr-bank.Start[layerNum], value)
}

func (B *BANK) Dump(startAddr uint16) {
	var val byte
	var line string
	var ascii string

	cpt := startAddr
	for j := 0; j < 16; j++ {
		fmt.Printf("%04X : ", cpt)
		line = ""
		ascii = ""
		for i := 0; i < 16; i++ {
			val = B.Read(cpt)
			if val != 0x00 && val != 0xFF {
				line = line + clog.CSprintf("white", "black", "%02X", val) + " "
			} else {
				line = fmt.Sprintf("%s%02X ", line, val)
			}
			if _, ok := trace.PETSCII[val]; ok {
				ascii += fmt.Sprintf("%s", string(trace.PETSCII[val]))
			} else {
				ascii += "."
			}
			cpt++
		}
		fmt.Printf("%s - %s\n", line, ascii)
	}
}

func (B *BANK) Show() {
	B.Layouts[*B.Selector&0x1F].Show()
}
