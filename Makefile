all: build

build:
	go build -o irda button.go main.go leds.go
run:
	./irda
