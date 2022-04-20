
all:
	make main

main:
	go build -o newC64 cmd/newC64/*

vic: cmd/vicTest/main.go
	go build -o vic cmd/vicTest/*

mem: cmd/mem/main.go mem/config.go
	go build -o TestMem cmd/mem/*
