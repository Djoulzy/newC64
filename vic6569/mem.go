package vic6569

import (
	"fmt"
	"newC64/clog"
	"newC64/trace"
)

const PAGE_DIVIDER = 12

type VicMem struct {
	Pages      []int
	Layers     [][]byte
	LayersName []string
	Start      []uint16
}

func (VM *VicMem) Init(nbLayers int, size uint16) *VicMem {
	VM.Layers = make([][]byte, nbLayers)
	VM.LayersName = make([]string, nbLayers)
	nbPages := int(size >> PAGE_DIVIDER)
	VM.Pages = make([]int, nbPages)
	return VM
}

func (VM *VicMem) Attach(name string, layerNum int, pageNum int, content []byte) {
	nbPages := len(content) >> PAGE_DIVIDER
	VM.LayersName[layerNum] = name
	VM.Layers[layerNum] = content
	for i := 0; i < nbPages; i++ {
		VM.Pages[pageNum+i] = layerNum
	}
}

func (VM *VicMem) Read(addr uint16) byte {
	page := int(addr >> PAGE_DIVIDER)
	pageStart := uint16(page << PAGE_DIVIDER)
	return VM.Layers[VM.Pages[page]][addr-pageStart]
}

func (VM *VicMem) Show() {
	clog.CPrintf("darkgray", "black", "%10s: ", " ")
	for p := range VM.Pages {
		clog.CPrintf("darkgray", "black", "%02d ", p)
	}
	fmt.Printf("\n")
	for l := range VM.Layers {
		clog.CPrintf("darkgray", "black", "%10s: ", VM.LayersName[l])
		for _, page := range VM.Pages {
			if page == l {
				clog.CPrintf("white", "black", " X ")
			} else {
				clog.CPrintf("darkgray", "black", " - ")
			}
		}
		fmt.Printf("\n")
	}
}

func (VM *VicMem) Dump(startAddr uint16) {
	var val byte
	var line string
	var ascii string

	cpt := startAddr
	for j := 0; j < 16; j++ {
		fmt.Printf("%04X : ", cpt)
		line = ""
		ascii = ""
		for i := 0; i < 16; i++ {
			val = VM.Read(cpt)
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
