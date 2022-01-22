package cia6526

type Keyboard struct {
	col byte
	row byte
}

var Keyb_A Keyboard = Keyboard{col: 0b00000010, row: 0b00000100}
var Keyb_B Keyboard = Keyboard{col: 0b00001000, row: 0b00010000}

var keyMap = map[uint]Keyboard{
	0: {col: 0b00000000, row: 0b00000000},
	97: Keyb_A,
	98: Keyb_B,
}
