# Requirements for OS

* https://github.com/orangepi-xunlong/wiringOP Install wiringpi
```
git clone https://github.com/orangepi-xunlong/wiringOP.git
cd wiringOP && ./build
```

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

# GPIO Pins to Leds in Linux

Open file `cat /boot/armbianEnv.txt` and see whats in the `overlay_prefix=` line, mine value is `sun8i-h3`.

## Install the required tools

```
# apt-get install device-tree-compiler
```

## Compile overlay and install it 

**Keep note that i'm using the same prefix as we seen in the `overlay_prefix=` before**

```
dtc -@ -Hepapr -I dts -O dtb -o /boot/dtb/overlay/sun8i-h3-gpio-led.dtbo dts/gpio-led-overlay.dts
```

Add overlay line to `/boot/armbianEnv.txt`

```
overlays=gpio-led
```

## Or you can use precompiled:

If you are using OrangePI Lite (first model) when you can use precompiled binary:
```
mv dts/sun8i-h3-gpio-led.dtbo /boot/dtb/overlay/
```

Now reboot and enjoy!
