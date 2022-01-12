package register

import "newC64/memory"

type REG struct {
	Val        byte
	Ram        *byte
	LastAccess *memory.Access
}

func (R *REG) Init(mem *memory.MEM, index uint16, defaultVal byte) {
	R.Ram = &mem.Val[index]
	R.Val = defaultVal
	R.LastAccess = &mem.LastAccess[index]
}

func (R *REG) IsMofied() bool {
	if *R.LastAccess == memory.WRITE && (R.Val != *R.Ram) {
		return true
	} else {
		*R.LastAccess = memory.NONE
		return false
	}
}

func (R *REG) Input() byte {
	return *R.Ram
}

func (R *REG) Output(newVal byte) {
	R.Val = newVal
	*R.Ram = newVal
	*R.LastAccess = memory.NONE
}

func (R *REG) Reset() {
	*R.Ram = R.Val
	*R.LastAccess = memory.NONE
}

// func (R *REG) CiaRegWrite(addr uint16, val byte, access memory.Access) {
// 	var i uint16
// 	for i = 0; i < 16; i++ {
// 		M.Val[addr+(16*i)] = val
// 		M.LastAccess[addr+(16*i)] = access
// 	}
// }
