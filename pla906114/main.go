package pla906114

import "fmt"

func (P *PLA) Connect(mem interface{}) {
	fmt.Printf("pla906114 - Init\n")
	P.ram = mem.(memory)
}

// Init :
func (P *PLA) Init() {
	P.ram.Init()
}

func (P *PLA) Clear() {

}

func (P *PLA) Read(addr uint16) byte {
	fmt.Printf("pla906114 - Read - %04X\n", addr)
	return P.ram.Read(addr)
}

func (P *PLA) Write(addr uint16, value byte) {
	P.ram.Write(addr, value)
}
