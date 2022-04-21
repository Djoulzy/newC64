package mem

import "io/ioutil"

func LoadROM(size int, file string) []byte {
	val := make([]byte, size)
	if len(file) > 0 {
		data, err := ioutil.ReadFile(file)
		if err != nil {
			panic(err)
		}
		if len(data) != size {
			panic("Bad ROM Size")
		}
		for i := 0; i < size; i++ {
			val[i] = byte(data[i])
		}
	}
	return val
}

func Clear(zone []byte) {
	cpt := 0
	fill := byte(0x00)
	for i := range zone {
		zone[i] = fill
		cpt++
		if cpt == 0x40 {
			fill = ^fill
			cpt = 0
		}
	}
}

func Fill(zone []byte, val byte) {
	for i := range zone {
		zone[i] = val
	}
}
