# goOpenrazerAnim

On Windows, Razer's software comes with a lot of custom modes that are not implememted in their hardware (e.g. color wheel). This simple tool tries to fix this by allowing users to play APNG files on their keyboards.

*This is still very much work in progress and breaking changes WILL happen*

## Getting started

A short introduction to using the tool

### Building

``` sh
git clone https://github.com/x0wllaar/goOpenrazerAnim
cd goOpenrazerAnim
go get .
go build -o ../animPlay .
```

### Preparing animations

_On my computer, sending a single frame takes around 55 ms, which corresponds to around 17fps. Your values might be different_

1. Find the video you'd like to play on your keyboard and download it. _The video from the description of https://www.youtube.com/watch?v=U8zSOY09Wtk works fine_

2. Rescale the video to 22x6 and convert it to APNG (the dimensions of the keyboard):

``` sh
ffmpeg -i Rainbowring.mp4 -filter:v "scale=22:6" -r 17 anim.apng
```

Note the -r 17.

You can use other filters to speed up and slow down the video (such as setpts)

### Find the folder of your keyboard

From the OpenRazer wiki ( https://github.com/openrazer/openrazer/wiki/Using-the-keyboard-driver ):

Example of the device path:

``` sh
/sys/bus/hid/drivers/razerkbd/<DEVICE ID>/
```

The might be multiple device id's in the razerkbd folder, you need the one that has matrix_custom_frame and matrix_effect_custom on the inside

### Run the program

``` sh
./animPlay -anim /path/to/anim.apng -kbpath /path/to/device/id
```
