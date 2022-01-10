package memory

import "io/ioutil"

type MEM struct {
	Size     int
	ReadOnly bool
	Val      []byte
	Written  []bool
}

func (M *MEM) load(filename string) {

	M.ReadOnly = true

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	if len(data) != M.Size {
		panic("Bad ROM Size")
	}
	for i := 0; i < M.Size; i++ {
		M.Val[i] = byte(data[i])
	}
}

func (M *MEM) Clear() {
	cpt := 0
	fill := byte(0x00)
	for i := range M.Val {
		M.Val[i] = fill
		cpt++
		if cpt == 0x40 {
			fill = ^fill
			cpt = 0
		}
	}
}

func (M *MEM) Init(size int, file string) {
	M.Size = size
	M.Val = make([]byte, size)
	M.Written = make([]bool, size)
	if len(file) > 0 {
		M.load(file)
		M.ReadOnly = true
	} else {
		M.Clear()
		M.ReadOnly = false
	}
}

func (M *MEM) GetView(start int, size int) *MEM {
	new := MEM{
		Size:     size,
		ReadOnly: M.ReadOnly,
		Val:      M.Val[start : start+size],
		Written:  M.Written[start : start+size],
	}
	return &new
}

func (M *MEM) VicRegWrite(addr uint16, val byte) {
	var i uint16
	for i=0; i<10; i++ {
		M.Val[addr + i*0x40] = val
	}
}

func (M *MEM) CiaRegWrite(addr uint16, val byte) {
	var i uint16
	for i = 0; i < 16; i++ {
		M.Val[addr + (16*i)] = val
	}
}