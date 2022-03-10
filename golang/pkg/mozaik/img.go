package mozaik

import (
	"fmt"
	"image"
	"image/color"
	_ "image/png"
	"os"
	"strings"

	"github.com/gosimple/slug"
	"github.com/nfnt/resize"
)

type LegoImage struct {
	image.Image
	Name             string
	Colors           Colors
	ColorList        ColorList
	ColorReplacement map[string]string
	Matrix           []Colors
	OriginalImage    image.Image
	Width            int
	Height           int
}

var colorCache map[color.Color]Color = map[color.Color]Color{}

func NewFromPath(path string, name string, colors Colors, width int, height int, replacement map[string]string) (*LegoImage, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	return NewLegoImage(img, name, colors, width, height, replacement)
}

func NewLegoImage(img image.Image, name string, colors Colors, width int, height int, replacement map[string]string) (*LegoImage, error) {
	var (
		legoImg LegoImage
		err     error
	)

	legoImg.ColorList = ColorList{}
	legoImg.Colors = colors
	legoImg.Image = img
	legoImg.Name = name
	legoImg.ColorReplacement = replacement
	legoImg.OriginalImage = img
	legoImg.Width = width
	legoImg.Height = height
	if err != nil {
		return nil, err
	}

	legoImg.Resize(15*width, 15*height)
	return &legoImg, nil
}

func (img *LegoImage) Resize(width int, height int) {
	img.Image = resize.Thumbnail(uint(width), uint(height), img.Image, resize.NearestNeighbor)
	img.matchColors()
}

func (img *LegoImage) matchColors() {
	img.Matrix = make([]Colors, img.Bounds().Dx())

	img.ColorList = ColorList{}
	matchedImage := image.NewRGBA(img.Image.Bounds())
	for x := 0; x < img.Bounds().Dx(); x++ {
		img.Matrix[x] = make(Colors, img.Bounds().Dy())
		for y := 0; y < img.Bounds().Dy(); y++ {
			colorAtXY := img.At(x, y)
			colorFound, ok := colorCache[colorAtXY]
			if !ok {
				colorFound = img.Colors.Find(colorAtXY)
				colorCache[colorAtXY] = colorFound
			}
			replacement, ok := img.ColorReplacement[colorFound.Name]
			if ok {
				colorReplacement, err := img.Colors.Get(replacement)
				if err == nil {
					colorFound = colorReplacement
				}
			}
			img.Matrix[x][y] = colorFound
			img.ColorList[colorFound.Name] = img.ColorList[colorFound.Name] + 1
			matchedImage.Set(x, y, color.RGBA{uint8(colorFound.R), uint8(colorFound.G), uint8(colorFound.B), 0xff})
		}
	}
	img.Image = matchedImage
}

func (img LegoImage) String() string {
	message := []string{
		"",
		fmt.Sprintf("Width: %d", img.Bounds().Dx()),
		fmt.Sprintf("Height: %d", img.Bounds().Dy()),
		"----",
		img.ColorList.String(),
	}
	return strings.Join(message, "\n")
}

func (img *LegoImage) Mozaik() *Mozaik {
	img.Resize(10*img.Width, 10*img.Height)
	img.Resize(5*img.Width, 5*img.Height)
	img.Resize(img.Width, img.Height)
	return &Mozaik{
		Name:          slug.Make(img.Name),
		Image:         img.Matrix,
		Width:         img.Width,
		Height:        img.Height,
		OriginalImage: img.OriginalImage,
	}
}

func (img LegoImage) ToMozaik() *Mozaik {
	mozaik := Mozaik{}
	return &mozaik
}
