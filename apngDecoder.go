package main

import (
	"image"
	"os"

	"github.com/kettek/apng"
	"github.com/sirupsen/logrus"
)

type APNGDecoder struct {
	filePath string
	frames   chan image.Image
	stop     chan bool
	running  bool
}

func NewAPNGDecoder(filePath string, frames chan image.Image) *APNGDecoder {
	return &APNGDecoder{
		filePath: filePath,
		frames:   frames,
		stop:     make(chan bool),
		running:  false,
	}
}

func (cd *APNGDecoder) Stop() {
	if cd.running {
		cd.stop <- true
	}
}

func (cd *APNGDecoder) Start() {
	if !cd.running {
		cd.running = true
		go cd.decodeAPNGFrames()
	}
}

func (cd *APNGDecoder) decodeAPNGFrames() {
	defer func() { cd.running = false }()

	logrus.Infof("Reading animation file from %s", cd.filePath)
	inputFile, err := os.Open(cd.filePath)
	if err != nil {
		logrus.Panicf("Error opening file: %v", err)
	}
	defer inputFile.Close()

	animation, err := apng.DecodeAll(inputFile)
	if err != nil {
		logrus.Panicf("Error decoding APNG %v", err)
	}

	logrus.Infof("Read %d frames", len(animation.Frames))

	for {
		for fn, frame := range animation.Frames {
			frameImage := frame.Image
			frameBounds := frameImage.Bounds()

			if (frameBounds.Max.X != EXPECTED_WIDTH) || (frameBounds.Max.Y != EXPECTED_HEIGHT) {
				logrus.Panicf(
					"Unexpected resolution in frame %d, expected %dx%d, got %dx%d",
					fn,
					EXPECTED_WIDTH, EXPECTED_HEIGHT,
					frameBounds.Max.X, frameBounds.Max.Y,
				)
			}

			select {
			case cd.frames <- frameImage:
			  	logrus.Tracef("Sent frame %d", fn)
			case stopMessage := <-cd.stop:
				if stopMessage {
		  			return
				}
			}
		}
	}
}
