* Check if kernel has exported leds in /sys/class/leds/ and do not use them in the app
* Read https://www.hpl.hp.com/personal/Jean_Tourrilhes/IrDA/IrNET.html 
* Move all source to src/ to clean the things up
* Code a ppp handler than will handle the connections and dropouts
* Integrate the access point possibility to the source from my other project (using hostpad)
* Code some config generation functions that will generate a required ppp configs with only minimal information passed.
* Integrate bonjour service from my other project into the code base
* Connect bahaviour from ppp to gpio led (blue) maybe we can make some noise with the traffic activity also
* Compile in options to disable gpio or other features (only for dev and testing)
* Gpio buttons behaviour, left button to restart ppp, right button safely shutdown the device, emergency wifi hotspot mode when pressing right button over 8 times.
* Write function to check the file permissions (+x execute) for every binary that are required for the project
