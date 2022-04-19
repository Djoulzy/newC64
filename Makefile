
all:
	make main

main:
	go build -o newC64 cmd/newC64/*

vic: cmd/vicTest/main.go
	go build -o vic cmd/vicTest/*
