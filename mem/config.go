package mem

import (
	"fmt"
	"io/ioutil"
	"newC64/clog"
	"newC64/trace"
)

const PAGE_DIVIDER = 12

type MEMAccess interface {
	MRead([]byte, uint16) byte
	MWrite([]byte, uint16, byte)
}

type CONFIG struct {
	Pages      []int
	Layers     [][]byte
	LayersName []string
	Start      []uint16
	Accessor   []MEMAccess
}

func Init(nbLayers int, size int) CONFIG {
	C := CONFIG{}
	C.Layers = make([][]byte, nbLayers)
	C.LayersName = make([]string, nbLayers)
	C.Start = make([]uint16, nbLayers)
	nbPages := int(size >> PAGE_DIVIDER)
	C.Pages = make([]int, nbPages)
	C.Accessor = make([]MEMAccess, nbLayers)
	return C
}

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

func (C *CONFIG) Attach(name string, layerNum int, pageNum int, start uint16, content []byte) {
	nbPages := len(content) >> PAGE_DIVIDER
	C.LayersName[layerNum] = name
	C.Layers[layerNum] = content
	C.Start[layerNum] = start
	for i := 0; i < nbPages; i++ {
		C.Pages[pageNum+i] = layerNum
	}
	C.Accessor[layerNum] = C
}

func (C *CONFIG) Accessors(layerNum int, access MEMAccess) {
	C.Accessor[layerNum] = access
}

func (C *CONFIG) Read(addr uint16) byte {
	layerNum := C.Pages[int(addr>>PAGE_DIVIDER)]
	// return C.Layers[layerNum][addr-C.Start[layerNum]]
	return C.Accessor[layerNum].MRead(C.Layers[layerNum], addr-C.Start[layerNum])
}

func (C *CONFIG) MRead(mem []byte, addr uint16) byte {
	return mem[addr]
}

func (C *CONFIG) MWrite(meme []byte, addr uint16, val byte) {

}

func (C *CONFIG) Show() {
	clog.CPrintf("dark_gray", "black", "\n%10s: ", "Pages")
	for p := range C.Pages {
		clog.CPrintf("dark_gray", "black", " %02d  ", p)
	}
	clog.CPrintf("dark_gray", "black", "\n%10s: ", "Start Addr")
	for p := range C.Pages {
		clog.CPrintf("light_gray", "black", "%04X ", p<<PAGE_DIVIDER)
	}
	fmt.Printf("\n")
	for l := range C.Layers {
		clog.CPrintf("light_gray", "black", "%10s: ", C.LayersName[l])
		for _, page := range C.Pages {
			if page == l {
				clog.CPrintf("black", "white", "     ")
			} else {
				clog.CPrintf("black", "dark_gray", "     ")
			}
		}
		fmt.Printf("\n")
	}
	fmt.Printf("\n")
}

func (C *CONFIG) Dump(startAddr uint16) {
	var val byte
	var line string
	var ascii string

	cpt := startAddr
	for j := 0; j < 16; j++ {
		fmt.Printf("%04X : ", cpt)
		line = ""
		ascii = ""
		for i := 0; i < 16; i++ {
			val = C.Read(cpt)
			if val != 0x00 && val != 0xFF {
				line = line + clog.CSprintf("white", "black", "%02X", val) + " "
			} else {
				line = fmt.Sprintf("%s%02X ", line, val)
			}
			if _, ok := trace.PETSCII[val]; ok {
				ascii += fmt.Sprintf("%s", string(trace.PETSCII[val]))
			} else {
				ascii += "."
			}
			cpt++
		}
		fmt.Printf("%s - %s\n", line, ascii)
	}
}
