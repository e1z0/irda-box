package main

import (
	. "github.com/cyoung/rpi"
	"time"
	"fmt"
	"flag"
        "os"
)

func init() {
	WiringPiSetup()
        initializeLeds()
}

func TestLeds() {
  for {
     for _, led := range leds {
       fmt.Printf("up led: %s\n",led.Color)
       led.On()
       time.Sleep(2*time.Second)
       fmt.Printf("down led: %s\n",led.Color)
       led.Off()
       time.Sleep(2*time.Second)
       time.Sleep(1*time.Second)
     }
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
	go GpioEvents()
	for {
		// just silly loop
		time.Sleep(10*time.Second)
	}
}
