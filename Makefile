all: build

build:
	go build -o irda button.go main.go leds.go commands.go web.go functions.go settings.go ppp.go upload.go
run:
	./irda
