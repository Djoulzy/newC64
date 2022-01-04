package main

type CPU interface {
	Init()
	Reset()
	NextCycle()
}

type MEM interface {
	Init()
	Clear()
	Read(uint16) byte
	Write(uint16, byte)
}
