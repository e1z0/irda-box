# Requirements for OS

* https://github.com/orangepi-xunlong/wiringOP Install wiringpi
```
git clone https://github.com/orangepi-xunlong/wiringOP.git
cd wiringOP && ./build
``

# GPIO Pins

Mentioned board pins means that the wire goes directly to the physical pin number as shown in the picture https://žn.lt/wiki/Vaizdas:OPiLite_pinout.jpg

UPS+ -> Board 5V (VCC)
UPS- -> Board GND

Lower button pin1 -> Board pin 5
Lower button pin2 -> Board pin 39

Upper button pin1 -> Board pin 7
Upper button pin2 -> Board pin 39

Upper diode (blue) + -> Board pin 15
Middle diode (green) + -> Board pin 19
Lower diode (red) + -> Board pin 21
All leds ground (the short leg) - goes to the ground
