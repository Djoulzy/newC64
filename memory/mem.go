package memory

import (
	"fmt"
	"io/ioutil"
	"log"
)

// Init :
func (M *MEM) Init() {
	M.Cells = make([]byte, M.Size)
	M.readOnly = false

	fmt.Printf("ram4164 - Init\n")
}

func (M *MEM) Load(filename string) {
	M.Cells = make([]byte, M.Size)
	M.readOnly = true

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	if len(data) != M.Size {
		panic("Bad ROM Size")
	}
	for i := 0; i < M.Size; i++ {
		M.Cells[i] = byte(data[i])
	}
}

func (M *MEM) Clear() {
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
}

func (M *MEM) Read(addr uint16) byte {
	return M.Cells[addr]
}

func (M *MEM) Write(addr uint16, value byte) {
	if M.readOnly {
		log.Fatal("Try to write protected area")
	}
	M.Cells[addr] = value
}

func (M *MEM) GetView(start int, size int) interface{} {
	view := MEM{
		Size:          size,
		readOnly:      M.readOnly,
		StartLocation: 0,
		Cells:         M.Cells[start : start+size-1],
	}
	return &view
}

func (M *MEM) Dump(startAddr uint16) {
	cpt := startAddr
	fmt.Printf("\n")
	for j := 0; j < 16; j++ {
		fmt.Printf("%04X : ", cpt)
		for i := 0; i < 16; i++ {
			fmt.Printf("%02X ", M.Cells[cpt])
			cpt++
		}
		fmt.Println()
	}
}
