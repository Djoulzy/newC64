package roms

import "io/ioutil"

func (R *ROM) Init(filename string, size int) {
	R.Cells = make([]byte, size)

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	if len(data) != size {
		panic("Bad ROM Size")
	}
	for i := 0; i < size; i++ {
		R.Cells[i] = byte(data[i])
	}
}

func (R *ROM) Read(addr uint16) byte {
	return 0
}
