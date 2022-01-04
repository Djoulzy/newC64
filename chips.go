package main

type MEM interface {
	Init()
	Clear()
	Read(uint16) byte
	Write(uint16, byte)
}

type PLA interface {
	Connect(interface{})
	Init()
	Clear()
	Read(uint16) byte
	Write(uint16, byte)
}

type CPU interface {
	Init(interface{})
	Reset()
	NextCycle()
}
