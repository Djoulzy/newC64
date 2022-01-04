package ram4164

import "fmt"

// Init :
func (m *RAM) Init() {
	cpt := 0
	fill := byte(0x00)
	for i := range m.Cells {
		m.Cells[i] = fill
		cpt++
		if cpt == 0x40 {
			fill = ^fill
			cpt = 0
		}
	}
	fmt.Printf("ram4164 - Init\n")
}

func (m *RAM) Clear() {

}

func (m *RAM) Read(addr uint16) byte {
	return m.Cells[addr]
}

func (m *RAM) Write(addr uint16, value byte) {
	m.Cells[addr] = value
}
