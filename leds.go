package main


import (
      . "github.com/cyoung/rpi"
)

var (
// uper blue led
UPPER_LED = 8
// middle green led
MIDDLE_LED = 11
// lower red led
LOWER_LED = 12
leds []Led
)

type Led struct {
Pin int
Name string
State bool
Color string
}

func (l *Led) On() {
 DigitalWrite(l.Pin,HIGH)
 l.State = true
}

func (l *Led) Off() {
 DigitalWrite(l.Pin,LOW)
 l.State = false
}

func (l *Led) Toggle() {
 if l.State {
    DigitalWrite(l.Pin,LOW)
    l.State = false
 } else {
    DigitalWrite(l.Pin,HIGH)
    l.State = true
 }
}


func initializeLeds() {
  red := Led{Pin: LOWER_LED, Name: "Critical Led", State: false, Color: "red"}
  blue := Led{Pin: MIDDLE_LED, Name: "System status Led", State: false, Color: "blue"}
  green := Led{Pin: UPPER_LED, Name: "IrDA status Led", State: false, Color: "green"}
  leds = append(leds,red,blue,green)
  for _,ld := range leds{
     PinMode(ld.Pin, OUTPUT)
  }
}
