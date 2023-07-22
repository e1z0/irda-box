package main

import (
      . "github.com/cyoung/rpi"
        "fmt"
        "time"
)

var (
SETUP_AVAILABLE=true
)

func GpioEvents() {
    PinMode(1, INPUT)
    PinMode(2, INPUT)
    btn_pushed_lower := 0
    btn_pushed_upper := 0
    last_time := time.Now().UnixNano() / 1000000
    // button lower
    go func() {
    for pinas := range WiringPiISR(1, INT_EDGE_FALLING) {
              if (pinas > -1 && SETUP_AVAILABLE) {
                 n := time.Now().UnixNano() / 1000000
                 delta := n - last_time
                 if delta > 1000 {
                 // reset counter
                 fmt.Printf("Reseting counter for lower button\n")
                 btn_pushed_lower=0
                 }
                 if delta > 300 { //software debouncing
                        fmt.Printf("Lower Button pressed: %d times\n",btn_pushed_lower)
                        if (btn_pushed_lower > 7 && SETUP_AVAILABLE) {
   //                     SETUP_AVAILABLE=false
                        //WifiSetup()
                        fmt.Printf("Setup sequences received from the lower button\n")
                        btn_pushed_lower=0
                        }
                        last_time = n
                        btn_pushed_lower++
                 }
               }
   }
   }()
   // button upper
    go func() {
    for pinas := range WiringPiISR(2, INT_EDGE_FALLING) {
              if (pinas > -1 && SETUP_AVAILABLE) {
                 n := time.Now().UnixNano() / 1000000
                 delta := n - last_time
                 if delta > 1000 {
                 // reset counter
                 fmt.Printf("Reseting counter for upper button\n")
                 btn_pushed_upper=0
                 }
                 if delta > 300 { //software debouncing
                        fmt.Printf("Upper Button pressed: %d times\n",btn_pushed_upper)
                        if (btn_pushed_upper > 7 && SETUP_AVAILABLE) {
 //                       SETUP_AVAILABLE=false
                        //WifiSetup()
                        fmt.Printf("Setup sequences received from the upper button\n")
                        btn_pushed_upper=0
                        }
                        last_time = n
                        btn_pushed_upper++
                 }
               }
   }
   }()
}
