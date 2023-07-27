# IrDA Box

<img src="https://raw.githubusercontent.com/e1z0/irda-box/master/pics/prototype1/irdabox-firstprototype_final_case_view2.jpeg" data-canonical-src="https://raw.githubusercontent.com/e1z0/irda-box/master/pics/prototype1/irdabox-firstprototype_final_case_view2.jpeg" width="600"/>

## Introduction

The Wifi-to-IrDA Internet Sharing Bridge (**IrDA Box**) is an innovative and cutting-edge device designed to address the connectivity gap between older devices with IrDA ports, such as legacy PDAs and Handheld PCs, and modern wireless networks. By harnessing the power of wireless optical communication, this bridge seamlessly connects devices with only IrDA support to the vast resources of the internet via a WiFi network. This groundbreaking technology ensures that even outdated gadgets can now access the internet, opening up new possibilities for productivity, communication, and entertainment.

## What is for ?

1. IrDA to WiFi Conversion: The device acts as a bridge, converting signals between Infrared Data Association (IrDA) and WiFi, effectively translating wireless signals to enable internet access for IrDA-supported devices.
2. Seamless Integration: The device seamlessly integrates into existing WiFi networks, making the setup process hassle-free and user-friendly. Users can easily connect their IrDA-equipped devices to the WiFi network, ensuring a smooth and uninterrupted internet experience.
3. Compact and Portable: The IrDA-to-WiFi bridge is designed to be compact and portable, making it easy to carry and use on the go. Its sleek design ensures it blends seamlessly into any environment, lol.
4. Compatibility: The bridge is compatible with a wide range of legacy devices that support the IrDA standard, including PDAs, Handheld PCs, and other similar gadgets. Now, these devices can join the modern internet age without requiring expensive upgrades or replacements.
Secure Connection: Security is of paramount importance, and the bridge is equipped with advanced encryption protocols to safeguard user data and ensure a secure internet connection.
5. Low Power Consumption: The device is engineered to consume minimal power, maximizing its operational efficiency and extending battery life for portable devices.

# Where to buy ?

I'm sorry at the time the device is in prototype stage and no ETA is defined.


# Build yourself

You can make the device yourself, just make sure you have these components, required for the device to be built:

* [OrangePI Lite](https://www.aliexpress.com/item/1005002557347741.html) or compatible OrangePI device with H3 SoC (Others such as RaspberryPI cant be used but currently not supported yet!
* [IrDA Adapter](https://www.ebay.com/itm/385652134546)
* Spare MicrosSD card with at least 8GB in size
* Write [Armbian 22.02.1 Focal](https://stpete-mirror.armbian.com/archive/orangepilite/archive/Armbian_22.02.1_Orangepilite_focal_current_5.15.25.img.xz) to MicroSD Card
* [Type-C 15W 3A 18650 Lithium Battery Charger  UPS Power Supply 5V](https://www.aliexpress.com/item/1005004863257598.html)
* 2x [buttons](https://www.aliexpress.com/item/32703664513.html)
* 3x [led diodes](https://www.aliexpress.com/item/1005005182451381.html)
* 1x 5mm [switch button](https://www.aliexpress.com/item/4001207529493.html)
* 3D Printer for [Case]() printing
* Some soldering skills
* [Soldering iron](https://www.aliexpress.com/item/1005005368783447.html)
* [Soldering wire](https://www.aliexpress.com/item/4001230482375.html)
* [Single Core Copper Wire](https://www.aliexpress.com/item/1005001918707461.html), wire can be different, it does not matter.

The build process of prototype 1 case is [shown here](prototype1.md)

## Requirements for OS

* Build tools

```
apt-get install git build-essential
```

* https://github.com/orangepi-xunlong/wiringOP Install wiringpi
```
git clone https://github.com/orangepi-xunlong/wiringOP.git
cd wiringOP && ./build
```

## GPIO Pins

Mentioned board pins means that the wire goes directly to the physical pin number as shown in the picture 

<img src="https://raw.githubusercontent.com/e1z0/irda-box/master/pics/irda_box_wiring.jpeg" data-canonical-src="https://raw.githubusercontent.com/e1z0/irda-box/master/pics/irda_box_wiring.jpeg" width="600"/>

| Component pin | Board pin |
| ------------- | ------------- |
| UPS/Battery **+**  | 5V (VCC)  |
| UPS/Battery **-**  | GND  |
| Left button pin1   | 5 |
| Left button pin2   | 39 |
| Right button pin1  | 7 |
| Right button pin2  | 39 |
| Upper diode (blue) + | 15 |
| Middle diode (green) + | 19 |
| Lower diode (red) + | 21 |

All leds grounds (the short leg) - goes to the ground

### GPIO Pins to Leds in Linux

Open file `cat /boot/armbianEnv.txt` and see whats in the `overlay_prefix=` line, mine value is `sun8i-h3`.

#### Install the required tools

```
# apt-get install device-tree-compiler
```

### Compile overlay and install it 

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
