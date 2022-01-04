package pla906114

import "fmt"

func (P *PLA) Init() {
	P.mem[RAM].Init()
	P.setting = 7
	fmt.Printf("pla906114 - Settings: %d\n", P.setting)
}

func (P *PLA) Clear() {
	P.mem[RAM].Clear()
}

func (P *PLA) Load(filename string) {

}

func (P *PLA) Attach(mem interface{}, memtype interface{}) {
	selectedType := memtype.(MemType)
	P.mem[selectedType] = mem.(memory)
}

func (P *PLA) getChip(addr uint16) MemType {
	switch {
	case addr < 0xA000:
		return RAM
	case addr < 0xC000:
		if P.setting&3 == 3 {
			return BASIC
		} else {
			return RAM
		}
	case addr < 0xD000:
		return RAM
	case addr < 0xE000:
		if P.setting&3 == 0 {
			return RAM
		} else if P.setting&4 == 0 {
			return CHAR
		} else if P.setting == 1 {
			return RAM
		} else {
			return IO
		}
	default:
		if P.setting&3 < 2 {
			return RAM
		} else {
			return KERNAL
		}
	}
}

func (P *PLA) Read(addr uint16) byte {
	dest := P.getChip(addr)
	fmt.Printf("pla906114 - Read - %04X - Zone: %d\n", addr, dest)
	return P.mem[dest].Read(addr)
}

func (P *PLA) Write(addr uint16, value byte) {
	dest := P.getChip(addr)
	P.mem[dest].Write(addr, value)
}

func (P *PLA) Dump(startAddr uint16) {
	cpt := startAddr
	fmt.Printf("\n")
	for j := 0; j < 16; j++ {
		fmt.Printf("%04X : ", cpt)
		for i := 0; i < 16; i++ {
			dest := P.getChip(cpt)
			fmt.Printf("%02X ", P.mem[dest].Read(cpt))
			cpt++
		}
		fmt.Println()
	}
}
