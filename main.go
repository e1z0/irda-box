package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	. "github.com/cyoung/rpi"
)

var (
	ENABLE_GPIO = "false"
)

func init() {
	PPPDaemon = PPP{PPPBinary: "./sh-test"}
	if ENABLE_GPIO == "true" {
		WiringPiSetup()
		initializeLeds()
	}
}

func TestLeds() {
	if ENABLE_GPIO == "true" {
		for {
			for _, led := range leds {
				fmt.Printf("up led: %s\n", led.Color)
				led.On()
				time.Sleep(2 * time.Second)
				fmt.Printf("down led: %s\n", led.Color)
				led.Off()
				time.Sleep(2 * time.Second)
				time.Sleep(1 * time.Second)
			}
		}
	} else {
		fmt.Printf("Gpio support is not compiled in\n")
	}
}

func main() {
	fmt.Printf("Launched!\n")
	testleds := flag.Bool("testleds", false, "Test leds only")
	flag.Parse()
	if *testleds {
		fmt.Printf("Leds tests mode\n")
		TestLeds()
		os.Exit(0)
	}
	//go GpioEvents()
	batteries := batteryInfo()
	for i, bat := range batteries {
		log.Printf("Battery %d: %#v\n", i, bat)
	}
	LoadProgramData()
	go StatusLoop()
	httpPool()
}
