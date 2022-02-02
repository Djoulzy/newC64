package memory

import (
	"fmt"
	"io/ioutil"
	"newC64/clog"
	"newC64/trace"
)

type Access byte

const (
	NONE Access = iota
	READ
	WRITE
)

type MEM struct {
	Size     int
	ReadOnly bool
	Val      []byte
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

func (M *MEM) Clear(pattern bool) {
	if pattern {
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
	} else {
		fill := byte(0xFF)
		for i := range M.Val {
			M.Val[i] = fill
		}
	}
}

func (M *MEM) Init(size int, file string) {
	M.Size = size
	M.Val = make([]byte, size)
	if len(file) > 0 {
		M.load(file)
		M.ReadOnly = true
	} else {
		M.Clear(false)
		M.ReadOnly = false
	}
}

func (M *MEM) GetView(start int, size int) *MEM {
	new := MEM{
		Size:     size,
		ReadOnly: M.ReadOnly,
		Val:      M.Val[start : start+size],
	}
	return &new
}

func (M *MEM) VicRegWrite(addr uint16, val byte, access Access) {
	var i uint16
	for i = 0; i < 10; i++ {
		M.Val[addr+i*0x40] = val
	}
}

func (M *MEM) CiaRegWrite(addr uint16, val byte, access Access) {
	var i uint16
	for i = 0; i < 16; i++ {
		M.Val[addr+(16*i)] = val
	}
}

func (M *MEM) Dump(startAddr uint16) {
	var val byte
	var line string
	var ascii string

	cpt := startAddr
	fmt.Printf("\n")
	for j := 0; j < 16; j++ {
		fmt.Printf("%04X : ", cpt)
		line = ""
		ascii = ""
		for i := 0; i < 16; i++ {
			val = M.Val[cpt]
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
