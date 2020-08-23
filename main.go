package main

import (
	"flag"
	"image"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
)

const EXPECTED_WIDTH = 22
const EXPECTED_HEIGHT = 6
const TRIGGER_STRING = "1"

var (
	keyboardPath *string = flag.String("kbpath", "/sys/bus/hid/drivers/razerkbd/0003:1532:0228.0006", "Path to the directory with the keyboard file")
	animFps      *int    = flag.Int("fps", -1, "Speed of the animation (will insert a delay of 1/fps seconds, -1 means no delay)")
	animFile     *string = flag.String("anim", "./anim.apng", "Path to the animation file")
)

func main() {
	logrus.SetLevel(logrus.InfoLevel)
	logrus.Infoln("Parsing args")
	flag.Parse()

	frames := make(chan image.Image)

	logrus.Infof("Starting frame producer")
	decoder := NewAPNGDecoder(*animFile, frames)
	decoder.Start()
	defer decoder.Stop()

	logrus.Infof("Starting frame sender")
	sender := NewFrameSender(*keyboardPath, frames, time.Duration(int64(time.Second)/int64(*animFps)))
	sender.Start()
	defer sender.Stop()

	mainSignal := make(chan os.Signal)
	signal.Notify(mainSignal, os.Interrupt, syscall.SIGHUP, syscall.SIGTERM)
	<-mainSignal
	logrus.Infoln("Stopping gracefully")

}
