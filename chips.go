package main

type MEM interface {
	Init()
	Clear()
	Load(string)
	Read(uint16) byte
	Write(uint16, byte)
	Dump(uint16)
}

type PLA interface {
	Init()
	Clear()
	Load(string)
	Attach(interface{}, interface{})
	Read(uint16) byte
	Write(uint16, byte)
	Dump(uint16)
}

type CPU interface {
	Init(interface{})
	Reset()
	NextCycle()
}
