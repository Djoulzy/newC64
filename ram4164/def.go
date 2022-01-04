package ram4164

const (
	memorySize  = 65536
	stackStart  = 0x0100
	stackEnd    = 0x01FF
	screenStart = 0x0400
	screenEnd   = 0x07FF
	charStart   = 0xD000
	charEnd     = 0xDFFF
	colorStart  = 0xD800
	colorEnd    = 0xDBFF
	intAddr     = 0xFFFA
	resetAddr   = 0xFFFC
	brkAddr     = 0xFFFE
	KernalStart = 0xE000
	KernalEnd   = 0xFFFF
	BasicStart  = 0xA000
	BasicEnd    = 0xC000
	vic2        = 0x4000
	vic3        = 0x8000
	vic4        = 0xC000
)

// RAM :
type RAM struct {
	Cells [memorySize]byte
}
