package memory

type MEM struct {
	Size          int
	readOnly      bool
	StartLocation int
	Cells         []byte
}
