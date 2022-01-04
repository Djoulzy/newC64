package main

import (
	"newC64/mos6510"
	"newC64/ram4164"
)

func main() {
	var cpu CPU
	var mem MEM

	mem = &ram4164.RAM{}
	mem.Init()

	cpu = &mos6510.CPU{}
	cpu.Init()
}
