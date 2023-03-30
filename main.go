package main

import (
	"flag"
	"fmt"
	"image"
	"image/png"
	"os"
)

func main() {
	// handle args
	confFilePath := flag.String("c", "", "Specify a path to a config file")
	canvasImage := flag.String("i", "", "Use an image as the canvas (overrides canvasColor)")
	extractFilePath := flag.String("e", "", "Extract colors from another file containing color codes")
	outFile := flag.String("o", "image.png", "Filename to use for output - should be a png")
	printConf := flag.Bool("p", false, "Print the config to stdout")
	wf := flag.Int("w", -1, "Set a canvasWidth (overrides the config file definition)")
	hf := flag.Int("h", -1, "Set a canvasHeight (overrides the config file definition)")
	xf := flag.Int("x", -1, "Set a yDistFactor (overrides the config file definition)")
	yf := flag.Int("y", -1, "Set an xDistFactor (overrides the config file definition)")
	flag.Parse()

	// read and load config
	var s *settings
	if *confFilePath != "" {
		var err error
		s, err = newSettingsFromConfig(*confFilePath)
		if err != nil {
			fmt.Printf("%s\n", err)
			os.Exit(1)
		}
		// set canvasWidth and canvasHeight for -a or if overridden
		if *wf != -1 {
			s.canvasWidth = *wf
		}
		if *hf != -1 {
			s.canvasHeight = *hf
		}
		if *xf != -1 {
			s.distX = s.canvasWidth / *xf
		}
		if *yf != -1 {
			s.distY = s.canvasHeight / *yf
		}
	} else if *extractFilePath != "" {
		var err error
		s, err = newSettingsFromFlags(*extractFilePath, *wf, *hf, *xf, *yf)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	} else {
		fmt.Println("Some source file is required - specifiy either '-c' or '-e'")
		os.Exit(1)
	}

	var dest *image.RGBA
	if *canvasImage != "" {
		imgFile, err := os.Open(*canvasImage)
		if err != nil {
			fmt.Println(err)
		}
		defer imgFile.Close()
		img, _, err := image.Decode(imgFile)
		dest = img.(*image.RGBA)
		if err != nil {
			fmt.Println(err)
		}
		s.fromSrcImage = true
	} else {
		// init canvas
		upLeft := image.Point{0, 0}
		lowRight := image.Point{s.canvasWidth, s.canvasHeight}

		dest = image.NewRGBA(image.Rectangle{upLeft, lowRight})
		drawBackground(dest, s.canvasColor)

	}

	drawManyRects(s, dest)

	// save to file
	out, err := os.Create(*outFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = png.Encode(out, dest)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if *printConf {
		printConfig(s)
	}
}
