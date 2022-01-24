package cia6526

type Keyboard struct {
	col byte
	row byte
}

const (
	PA0 byte = 0b00000001
	PA1 byte = 0b00000010
	PA2 byte = 0b00000100
	PA3 byte = 0b00001000
	PA4 byte = 0b00010000
	PA5 byte = 0b00100000
	PA6 byte = 0b01000000
	PA7 byte = 0b10000000
)

const (
	PB0 byte = 0b00000001
	PB1 byte = 0b00000010
	PB2 byte = 0b00000100
	PB3 byte = 0b00001000
	PB4 byte = 0b00010000
	PB5 byte = 0b00100000
	PB6 byte = 0b01000000
	PB7 byte = 0b10000000
)

var (
	Keyb_NULL       Keyboard = Keyboard{col: 0b00000000, row: 0b00000000}
	Keyb_0          Keyboard = Keyboard{col: PA4, row: PB3}
	Keyb_1          Keyboard = Keyboard{col: PA7, row: PB0}
	Keyb_2          Keyboard = Keyboard{col: PA7, row: PB3}
	Keyb_3          Keyboard = Keyboard{col: PA1, row: PB0}
	Keyb_4          Keyboard = Keyboard{col: PA1, row: PB3}
	Keyb_5          Keyboard = Keyboard{col: PA2, row: PB0}
	Keyb_6          Keyboard = Keyboard{col: PA2, row: PB3}
	Keyb_7          Keyboard = Keyboard{col: PA3, row: PB0}
	Keyb_8          Keyboard = Keyboard{col: PA3, row: PB3}
	Keyb_9          Keyboard = Keyboard{col: PA4, row: PB0}
	Keyb_A          Keyboard = Keyboard{col: PA1, row: PB2}
	Keyb_B          Keyboard = Keyboard{col: PA3, row: PB4}
	Keyb_C          Keyboard = Keyboard{col: PA2, row: PB4}
	Keyb_D          Keyboard = Keyboard{col: PA2, row: PB2}
	Keyb_E          Keyboard = Keyboard{col: PA1, row: PB6}
	Keyb_F          Keyboard = Keyboard{col: PA2, row: PB5}
	Keyb_G          Keyboard = Keyboard{col: PA3, row: PB2}
	Keyb_H          Keyboard = Keyboard{col: PA3, row: PB5}
	Keyb_I          Keyboard = Keyboard{col: PA4, row: PB1}
	Keyb_J          Keyboard = Keyboard{col: PA4, row: PB2}
	Keyb_K          Keyboard = Keyboard{col: PA4, row: PB5}
	Keyb_L          Keyboard = Keyboard{col: PA5, row: PB2}
	Keyb_M          Keyboard = Keyboard{col: PA4, row: PB4}
	Keyb_N          Keyboard = Keyboard{col: PA4, row: PB7}
	Keyb_O          Keyboard = Keyboard{col: PA4, row: PB6}
	Keyb_P          Keyboard = Keyboard{col: PA5, row: PB1}
	Keyb_Q          Keyboard = Keyboard{col: PA7, row: PB6}
	Keyb_R          Keyboard = Keyboard{col: PA2, row: PB1}
	Keyb_S          Keyboard = Keyboard{col: PA1, row: PB5}
	Keyb_T          Keyboard = Keyboard{col: PA2, row: PB6}
	Keyb_U          Keyboard = Keyboard{col: PA3, row: PB6}
	Keyb_V          Keyboard = Keyboard{col: PA3, row: PB7}
	Keyb_W          Keyboard = Keyboard{col: PA1, row: PB1}
	Keyb_X          Keyboard = Keyboard{col: PA2, row: PB7}
	Keyb_Y          Keyboard = Keyboard{col: PA3, row: PB1}
	Keyb_Z          Keyboard = Keyboard{col: PA1, row: PB4}
	Keyb_RETURN     Keyboard = Keyboard{col: PA0, row: PB1}
	Keyb_SPACE      Keyboard = Keyboard{col: PA7, row: PB4}
	Keyb_LSHIFT     Keyboard = Keyboard{col: PA1, row: PB7}
	Keyb_RSHIFT     Keyboard = Keyboard{col: PA6, row: PB4}
	Keyb_CRSR_DOWN  Keyboard = Keyboard{col: PA0, row: PB7}
	Keyb_CRSR_RIGHT Keyboard = Keyboard{col: PA0, row: PB2}
)

var keyMap = map[uint]Keyboard{
	0:          Keyb_NULL,
	13:         Keyb_RETURN,
	32:         Keyb_SPACE,
	48:         Keyb_0,
	49:         Keyb_1,
	50:         Keyb_2,
	51:         Keyb_3,
	52:         Keyb_4,
	53:         Keyb_5,
	54:         Keyb_6,
	55:         Keyb_7,
	56:         Keyb_8,
	57:         Keyb_9,
	97:         Keyb_A,
	98:         Keyb_B,
	99:         Keyb_C,
	100:        Keyb_D,
	101:        Keyb_E,
	102:        Keyb_F,
	103:        Keyb_G,
	104:        Keyb_H,
	105:        Keyb_I,
	106:        Keyb_J,
	107:        Keyb_K,
	108:        Keyb_L,
	109:        Keyb_M,
	110:        Keyb_N,
	111:        Keyb_O,
	112:        Keyb_P,
	113:        Keyb_Q,
	114:        Keyb_R,
	115:        Keyb_S,
	116:        Keyb_T,
	117:        Keyb_U,
	118:        Keyb_V,
	119:        Keyb_W,
	120:        Keyb_X,
	121:        Keyb_Y,
	122:        Keyb_Z,
	1073742049: Keyb_LSHIFT,
	1073742053: Keyb_RSHIFT,
	1073741905: Keyb_CRSR_DOWN,
	1073741903: Keyb_CRSR_RIGHT,
}