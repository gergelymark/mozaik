package mozaik

import (
	"encoding/json"
	"errors"
	"fmt"
	"image/color"
	"io/ioutil"
	"math"
	"os"
	"strings"

	"github.com/lucasb-eyer/go-colorful"
)

type Color struct {
	Name string `json:"name"`
	R    int    `json:"r"`
	G    int    `json:"g"`
	B    int    `json:"b"`
}

type Colors []Color

type ColorList map[string]int

type LegoColors struct {
	Colors []Color `json:"colors"`
}

func (color Color) AsColorful() colorful.Color {
	return colorful.Color{
		R: float64(color.R) / 255.0,
		G: float64(color.G) / 255.0,
		B: float64(color.B) / 255.0,
	}
}

func (baseColor Color) Dist(color color.Color) float64 {
	r, g, b, _ := color.RGBA()
	return baseColor.AsColorful().DistanceCIEDE2000(colorful.Color{
		R: float64(uint8(r)) / 255.0,
		G: float64(uint8(g)) / 255.0,
		B: float64(uint8(b)) / 255.0,
	})
}

func (colors Colors) Get(name string) (Color, error) {
	for _, color := range colors {
		if color.Name == name {
			return color, nil
		}
	}

	return Color{}, errors.New("color nor found")
}

func (colors Colors) Find(color color.Color) Color {
	var (
		nearestColor Color
		minDist      float64 = math.Inf(1)
	)

	for _, c := range colors {
		dist := c.Dist(color)
		if dist < minDist {
			minDist = dist
			nearestColor = c
		}
	}

	return nearestColor
}

func LoadLegoColors(path string) (Colors, error) {
	var err error
	jsonFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return nil, err
	}

	colorMap := LegoColors{}
	err = json.Unmarshal([]byte(byteValue), &colorMap)
	if err != nil {
		return nil, err
	}

	return colorMap.Colors, nil
}

func (colorList ColorList) String() string {
	colors := []string{}
	for colorName, count := range colorList {
		colors = append(colors, fmt.Sprintf("%s: %d", colorName, count))
	}

	return fmt.Sprintf("%s\n", strings.Join(colors, "\n"))
}
