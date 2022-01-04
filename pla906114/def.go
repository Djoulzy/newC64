package pla906114

type memory interface {
	Init()
	Clear()
	Read(uint16) byte
	Write(uint16, byte)
}

// RAM :
type PLA struct {
	ram        memory
}
