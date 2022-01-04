package pla906114

import "fmt"

func (P *PLA) Init() {
	fmt.Printf("pla906114 - Init\n")
	P.ram.Init()
}

func (P *PLA) Clear() {
	P.ram.Clear()
}

func (P *PLA) Attach(mem interface{}, memtype interface{}) {
	selectedType := memtype.(MemType)
	switch selectedType {
	case RAM:
		P.ram = mem.(memory)
	case KERNAL:
		fallthrough
	case BASIC:
		fallthrough
	case CHAR:
		P.rom[selectedType] = mem.(rom)
	}
}

func (P *PLA) Read(addr uint16) byte {
	fmt.Printf("pla906114 - Read - %04X\n", addr)
	return P.ram.Read(addr)
}

func (P *PLA) Write(addr uint16, value byte) {
	P.ram.Write(addr, value)
}
