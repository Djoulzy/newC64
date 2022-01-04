package ram4164

import "fmt"

// Init :
func (M *RAM) Init() {
	cpt := 0
	fill := byte(0x00)
	for i := range M.Cells {
		M.Cells[i] = fill
		cpt++
		if cpt == 0x40 {
			fill = ^fill
			cpt = 0
		}
	}
	fmt.Printf("ram4164 - Init\n")
}

func (M *RAM) Clear() {

}

func (M *RAM) Read(addr uint16) byte {
	fmt.Printf("ram4164 - Read - %04X\n", addr)
	return M.Cells[addr]
}

func (M *RAM) Write(addr uint16, value byte) {
	M.Cells[addr] = value
}
