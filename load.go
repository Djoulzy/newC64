package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
)

func LoadHex(mem MEM, code string) (uint16, error) {
	data := strings.Fields(code)
	tmp, _ := strconv.ParseUint(strings.TrimSuffix(data[0], ":"), 16, 16)
	start := uint16(tmp)
	fmt.Printf("Start: %04X\n", start)

	for i, val := range data[1:] {
		index := start + uint16(i)
		value, _ := strconv.ParseUint(val, 16, 8)
		mem.Write(index, byte(value))
	}
	return start, nil
}

func LoadFile(mem MEM, file string) (uint16, error) {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}
	text := string(content)
	return LoadHex(mem, text)
}

func LoadPRG(mem MEM, file string) (uint16, error) {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}
	startLoadMem := uint16(content[1]) << 8
	startLoadMem |= uint16(content[0])

	prgStart := (int(content[7])-0x30)*1000 + (int(content[8])-0x30)*100 + (int(content[9])-0x30)*10 + (int(content[10]) - 0x30)
	fmt.Printf("PRG: %04X\n", prgStart)
	for i, val := range content[2:] {
		mem.Write(startLoadMem+uint16(i), val)
	}
	return uint16(prgStart), nil
}
