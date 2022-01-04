package main

type MEM interface {
	Init()
	Clear()
	Read(uint16) byte
	Write(uint16, byte)
}

type ROM interface {
	Init(string, int)
	Read(uint16) byte
}

type PLA interface {
	Init()
	Clear()
	Attach(interface{}, interface{})
	Read(uint16) byte
	Write(uint16, byte)
}

type CPU interface {
	Init(interface{})
	Reset()
	NextCycle()
}
