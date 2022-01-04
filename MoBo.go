package main

import (
	"newC64/mos6510"
	"newC64/pla906114"
	"newC64/ram4164"
)

func main() {
	var cpu CPU
	var mem MEM
	var pla PLA

	mem = &ram4164.RAM{}

	pla = &pla906114.PLA{}
	pla.Connect(mem)
	pla.Init()

	cpu = &mos6510.CPU{}
	cpu.Init(pla)
}
