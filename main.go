package main

import (
	"flag"
	//	"fmt"
	"log"
	"os"
	"time"

	. "github.com/cyoung/rpi"
)

var (
	ENABLE_GPIO = "true"
        REQUIRED_BINS = []string{"ps","w","irdadump","pppd","lsmod","modprobe"}
)

func init() {
        CheckRequiredBins()
        pppbin,err := ReturnBinPath("pppd")
        if err != nil {
            log.Printf("Unable to find ppp binary, make sure you have install all required tools to run this application\n")
            log.Printf("Type ./irda -install to install all required software\n")
        }
		PPPDaemon = PPP{PPPBinary: pppbin, PPPArgs: []string{"call", "irda"}}
		IrUp = IRUpload{Disabled: true}
        // testing
//        PPPDaemon.PPPBinary = "./sh-test"
        ircomms, err := ReturnIrcommIfaces()
        if err != nil || len(ircomms) == 0 {
          log.Printf("No IrCOMM interfaces where found, the ircomm-tty kernel module is not loaded!\n")
          ok, err := ModProbe("ircomm-tty")
          if ok && err == nil {
              log.Printf("ircomm-tty kernel module have been loaded\n")
          } else {
              log.Printf("Unable to load ircomm-tty module: %s\n",err)
          }
        }
        ircomms,_ = ReturnIrcommIfaces()
        if len(ircomms) > 0 {
          log.Printf("IrCOMM interfaces: %d\n",len(ircomms))
        }

	if ENABLE_GPIO == "true" {
		WiringPiSetup()
		initializeLeds()
	}
}

func TestLeds() {
	if ENABLE_GPIO == "true" {
		for {
			for _, led := range leds {
				log.Printf("up led: %s\n", led.Color)
				led.On()
				time.Sleep(2 * time.Second)
				log.Printf("down led: %s\n", led.Color)
				led.Off()
				time.Sleep(2 * time.Second)
				time.Sleep(1 * time.Second)
			}
		}
	} else {
		log.Printf("Gpio support is not compiled in\n")
	}
}

func main() {
	log.Printf("Program launched!\n")
	testleds := flag.Bool("testleds", false, "Test leds only")
	flag.Parse()
	if *testleds {
		log.Printf("Leds tests mode\n")
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
