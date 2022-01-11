package register

type REG struct {
	val byte
	ram *byte
}

func (R *REG) Init(ramLoc *byte, defaultVal byte) {
	R.ram = ramLoc
	R.val = defaultVal
}
