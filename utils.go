package main

import (
	"bufio"
	"errors"
	"fmt"
	"image/color"
	"log"
	"math/rand"
	"os"
	"regexp"
	"strconv"
)

func randColor(colors []color.NRGBA) color.NRGBA {
	return colors[rand.Intn(len(colors))]
}

func hexToNRGBA(hex string) (clr color.NRGBA, err error) {
	var i int = 1
	if hex[0:1] != "#" {
		return color.NRGBA{0, 0, 0, 255}, errors.New("hex color codes should start with '#'")
	}
	if len(hex) != 7 {
		return color.NRGBA{0, 0, 0, 255}, errors.New("hex color codes should be 7 characters long")
	}

	r, _ := strconv.ParseUint(hex[i:i+2], 16, 8)
	g, _ := strconv.ParseUint(hex[i+2:i+4], 16, 8)
	b, _ := strconv.ParseUint(hex[i+4:i+6], 16, 8)

	a := 255 // no transparency
	ret := color.NRGBA{uint8(r), uint8(g), uint8(b), uint8(a)}
	//fmt.Println(ret)
	return ret, nil
}

func nrgbaToHex(clr color.NRGBA) (hex string) {
	r := strconv.FormatInt(int64(clr.R), 16)
	g := strconv.FormatInt(int64(clr.G), 16)
	b := strconv.FormatInt(int64(clr.B), 16)
	hex = fmt.Sprintf("#%s%s%s", r, g, b)
	return
}

func extractColors(path string) []string {
	var colors []string
	r := regexp.MustCompile("#?[a-fA-F0-9]{6}|#?[a-fA-F-9]{8}")

	//fmt.Println(path)
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		matches := r.FindAllString(scanner.Text(), -1)
		colors = append(colors, matches...)

		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
	}
	return colors
}

func randRange(min, max int) int {
	if min == max {
		return max
	} else {
		return rand.Intn(max-min) + min
	}
}
