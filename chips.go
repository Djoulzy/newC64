package main

type MEM interface {
	Init()
	Clear()
	Load(string)
	Read(uint16) byte
	Write(uint16, byte)
	GetView(int, int) interface{}
	Dump(uint16)
}

type PLA interface {
	Init()
	Clear()
	Load(string)
	Attach(interface{}, interface{}, int)
	Read(uint16) byte
	Write(uint16, byte)
	GetView(int, int) interface{}
	Dump(uint16)
}

type CPU interface {
	Init(interface{})
	Reset()
	NextCycle()
}
