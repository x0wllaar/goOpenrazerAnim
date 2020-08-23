package main

import (
	"fmt"
	"image"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

type FrameSender struct {
	keyboardDir string
	frames      chan image.Image
	stop        chan bool
	frameDelay  time.Duration
	running     bool
}

func NewFrameSender(keyboardDir string, frames chan image.Image, frameDelay time.Duration) *FrameSender {
	return &FrameSender{
		keyboardDir: keyboardDir,
		frames:      frames,
		stop:        make(chan bool),
		frameDelay:  frameDelay,
		running:     false,
	}
}

func (cfs *FrameSender) Stop() {
	if cfs.running {
		cfs.stop <- true
	}
}

func (cfs *FrameSender) Start() {
	if !cfs.running {
		cfs.running = true
		go cfs.sendFrames()
	}
}

func (cfs *FrameSender) sendFrames() {
	defer func() { cfs.running = false }()

	logrus.Infof("Opening driver files")
	customFrameFile, err := os.OpenFile(fmt.Sprintf("%s/matrix_custom_frame", cfs.keyboardDir), os.O_WRONLY, 0220)
	if err != nil {
		logrus.Panicf("Error opening driver file: %v", err)
	}
	defer customFrameFile.Close()

	customEffectFile, err := os.OpenFile(fmt.Sprintf("%s/matrix_effect_custom", cfs.keyboardDir), os.O_WRONLY, 0220)
	if err != nil {
		logrus.Panicf("Error opening driver file: %v", err)
	}
	defer customEffectFile.Close()

	for {
		select {
		case frame := <-cfs.frames:
			bTime := time.Now()
			frameSize := frame.Bounds().Max

			lineLen := (frameSize.X)*3 + 3
			logrus.Tracef("Line length %d bytes", lineLen)

			frameBytes := make([]byte, 0)
			for lineN := 0; lineN < frameSize.Y; lineN++ {
				lineBytes := make([]byte, lineLen)
				lineBytes[0] = byte(lineN)
				lineBytes[1] = byte(0)
				lineBytes[2] = byte(frameSize.X - 1)
				for colN := 0; colN < frameSize.X; colN++ {
					pixelColorR, pixelColorG, pixelColorB, _ := frame.At(colN, lineN).RGBA()
					baseOffset := 3 + (colN * 3)
					lineBytes[baseOffset+0] = byte(pixelColorR)
					lineBytes[baseOffset+1] = byte(pixelColorG)
					lineBytes[baseOffset+2] = byte(pixelColorB)
				}
				logrus.Tracef("Bytes %v", lineBytes)
				frameBytes = append(frameBytes, lineBytes...)
			}
			
			written, err := customFrameFile.Write(frameBytes)
			if err != nil {
				logrus.Errorf("Error writing frame %v", err)
			}
			if written != len(frameBytes) {
				logrus.Warnf("Failed to write bytes, %d written, %d expected", written, len(frameBytes))
			}
			logrus.Debugln("Frame written")

			_, err = customEffectFile.WriteString(TRIGGER_STRING)
			if err != nil {
				logrus.Errorf("Error triggering frame %v", err)
			}
			logrus.Debugln("Frame triggered")

			aTime := time.Now()
			logrus.Debugf("Frame time: %v", aTime.Sub(bTime))

			if cfs.frameDelay > 0 {
				time.Sleep(cfs.frameDelay)
			}

		case stopMsg := <-cfs.stop:
			if stopMsg {
				return
			}
		}
	}
}
