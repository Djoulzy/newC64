package main

import (
	"newC64/pla906114"
	"os"
)

func DumpMem(mem *pla906114.PLA, file string) error {
	var tmp []byte
	f, err := os.Create(file)
	if err != nil {
		return err
	}

	tmp = make([]byte, ramSize)
	for i := 0; i < ramSize; i++ {
		// tmp[i] = mem.Read(uint16(i))
		tmp[i] = mem.Mem[0].Val[i]
	}
	f.Write(tmp)

	return nil
}
