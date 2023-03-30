package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"image/color"
	"os"
	"strings"
)

const (
	defaultCanvasWidth  = 1000
	defaultCanvasHeight = 1000
	randRangeMin        = 2
	randRangeMax        = 40
	defaultJitterMax    = 100
)

type configFile struct {
	RectColors     []string `json:"rectColors"`
	CanvasWidth    int      `json:"canvasWidth"`
	CanvasHeight   int      `json:"canvasHeight"`
	CanvasColor    string   `json:"canvasColor"`
	XDistFactor    int      `json:"xDistFactor"`
	YDistFactor    int      `json:"yDistFactor"`
	JitterX        int      `json:"jitterX"`
	JitterY        int      `json:"jitterY"`
	JitterWidth    int      `json:"jitterWidth"`
	JitterHeight   int      `json:"jitterHeight"`
	RectWidth      int      `json:"rectWidth"`
	RectHeight     int      `json:"rectHeight"`
	PreserveSquare bool     `json:"preserveSquare"`
	Weight         int      `json:"weight,omitempty"`
	WeightedColor  string   `json:"weightedColor,omitempty"`
}

type settings struct {
	rectColors     []color.NRGBA
	canvasColor    color.NRGBA
	canvasWidth    int
	canvasHeight   int
	distX          int
	distY          int
	jitterX        int
	jitterY        int
	rectWidth      int
	rectHeight     int
	jitterWidth    int
	jitterHeight   int
	preserveSquare bool
	fromSrcImage   bool
}

func newSettingsFromFlags(path string, wf, hf, xf, yf int) (*settings, error) {
	clrs := extractColors(path)
	clrsNRGBA := make([]color.NRGBA, len(clrs))
	for i := range clrs {
		clrsNRGBA[i], _ = hexToNRGBA(clrs[i])
	}
	cc := randColor(clrsNRGBA)
	if wf == -1 {
		wf = defaultCanvasHeight
	}
	if hf == -1 {
		hf = defaultCanvasWidth
	}
	if xf == -1 {
		xf = randRange(randRangeMin, randRangeMax/2)
	}
	if yf == -1 {
		yf = randRange(randRangeMin, randRangeMax/2)
	}
	// These random ranges are arbitrary and dont always come out nicely
	// would be nice to figure out ranges/methods that produce nicer images.
	s := settings{
		rectColors:   clrsNRGBA,
		canvasColor:  cc,
		canvasWidth:  wf,
		canvasHeight: hf,
		distX:        wf / xf,
		distY:        hf / yf,
		jitterX:      0,
		jitterY:      0,
		rectWidth:    randRange(randRangeMin, randRangeMax),
		rectHeight:   randRange(randRangeMin, randRangeMax),
		jitterWidth:  0,
		jitterHeight: 0,
		fromSrcImage: false, // default to false
	}
	switch randRange(0, 3) {
	case 0:
		s.jitterWidth = randRange(0, defaultJitterMax)
	case 1:
		s.jitterHeight = randRange(0, defaultJitterMax)
	case 2:
		s.jitterX = randRange(0, defaultJitterMax)
	case 3:
		s.jitterY = randRange(0, defaultJitterMax)
	}

	return &s, nil
}

func newSettingsFromConfig(path string) (*settings, error) {
	cfile, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var c configFile
	err = json.Unmarshal(cfile, &c)
	if err != nil {
		return nil, err
	}
	// check configuration validity

	err = validateConfig(&c)
	if err != nil {
		return nil, err
	}

	// convert colors to NRGBA
	if c.Weight != 0 {
		c.RectColors = addWeight(c.RectColors, c.Weight, c.WeightedColor)
	}
	clrsNRGBA := make([]color.NRGBA, len(c.RectColors))
	for i := range c.RectColors {
		clrsNRGBA[i], _ = hexToNRGBA(c.RectColors[i])
	}

	cc, _ := hexToNRGBA(c.CanvasColor)
	s := settings{
		rectColors:     clrsNRGBA,
		canvasColor:    cc,
		canvasWidth:    c.CanvasWidth,
		canvasHeight:   c.CanvasHeight,
		distX:          c.CanvasWidth / c.XDistFactor,
		distY:          c.CanvasHeight / c.YDistFactor,
		jitterX:        c.JitterX,
		jitterY:        c.JitterY,
		rectWidth:      c.RectWidth,
		rectHeight:     c.RectHeight,
		jitterWidth:    c.JitterWidth,
		jitterHeight:   c.JitterHeight,
		preserveSquare: c.PreserveSquare,
		fromSrcImage:   false, // default to false
	}

	return &s, nil
}

func validateConfig(c *configFile) error {
	// validate colors
	for _, v := range c.RectColors {
		if v[0:1] != "#" || len(v) != 7 {
			return errors.New("invalid config: color codes should be hex " +
				"with 7 or 9 characters and start with '#'")
		}
	}

	reqProps := []string{
		"rectColors",
		"canvasWidth",
		"canvasHeight",
		"canvasColor",
		"xDistFactor",
		"yDistFactor",
		"rectWidth",
		"rectHeight",
	}

	//clumsy
	if c.RectColors == nil || len(c.RectColors) == 0 || c.CanvasWidth == 0 ||
		c.CanvasHeight == 0 || c.CanvasColor == "" || c.XDistFactor == 0 ||
		c.YDistFactor == 0 || c.RectWidth == 0 || c.RectHeight == 0 {
		fmt.Println(c.CanvasColor)
		return errors.New("invalid config: required properties [" +
			strings.Join(reqProps, ", ") + "] must be non-nil")
	}

	// handle mutually required props
	if (c.Weight != 0 && c.WeightedColor == "") ||
		(c.Weight == 0 && c.WeightedColor != "") {
		return errors.New("invalid config: if either 'weight' or " +
			"'weightedColor' is specified, the other is also required")
	}

	return nil
}

func addWeight(colors []string, weight int, hex string) []string {
	for i := 0; i < weight; i++ {
		colors = append(colors, hex)
	}
	return colors
}

func printConfig(s *settings) {
	var clrs []string = make([]string, len(s.rectColors))
	for i := range s.rectColors {
		clrs[i] = nrgbaToHex(s.rectColors[i])
	}
	var c configFile = configFile{
		RectColors:     clrs,
		CanvasWidth:    s.canvasWidth,
		CanvasHeight:   s.canvasHeight,
		CanvasColor:    nrgbaToHex(s.canvasColor),
		XDistFactor:    s.canvasWidth / s.distX,
		YDistFactor:    s.canvasHeight / s.distY,
		JitterX:        s.jitterX,
		JitterY:        s.jitterY,
		JitterWidth:    s.jitterWidth,
		JitterHeight:   s.jitterHeight,
		RectWidth:      s.rectWidth,
		RectHeight:     s.rectHeight,
		PreserveSquare: s.preserveSquare,
	}

	b, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(b))
}
