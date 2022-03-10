package mozaik

import (
	"bytes"
	"encoding/json"
	"image"
	"image/png"
	"io/ioutil"
	"os"
	"path"

	"github.com/gergelymark/mozaik/pkg/config"
	"github.com/vincent-petithory/dataurl"
)

type MozaikData struct {
	Name    string
	DataUrl string
	Width   int
	Height  int
}

type Mozaik struct {
	Name          string
	Image         []Colors
	Width         int
	Height        int
	OriginalImage image.Image `json:"-"`
}

type Mozaiks []Mozaik

func (mozaikData MozaikData) GetSourceImage() (image.Image, error) {
	dataURL, err := dataurl.DecodeString(mozaikData.DataUrl)
	if err != nil {
		return nil, err
	}
	sourceImage, _, err := image.Decode(bytes.NewReader(dataURL.Data))
	return sourceImage, err
}

func (mozaikData MozaikData) ToMozaik(legoColors Colors) (*Mozaik, error) {
	sourceImage, err := mozaikData.GetSourceImage()
	if err != nil {
		return nil, err
	}

	legoImg, err := NewLegoImage(sourceImage, mozaikData.Name, legoColors, mozaikData.Width, mozaikData.Height, map[string]string{})
	if err != nil {
		return nil, err
	}

	return legoImg.Mozaik(), nil
}

func (mozaik Mozaik) ToJSON() ([]byte, error) {
	return json.Marshal(mozaik)
}

func Load(name string) (*Mozaik, error) {
	var mozaik Mozaik
	mozaikPath := path.Join(config.Config.BasePath, name)
	jsonFile, err := os.Open(mozaikPath)
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()
	jsonBytes, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(jsonBytes, &mozaik)
	if err != nil {
		return nil, err
	}

	return &mozaik, nil
}

func (mozaik Mozaik) Save() error {
	mozaikPath := path.Join(config.Config.BasePath, mozaik.Name)
	err := os.MkdirAll(mozaikPath, 0755)
	if err != nil {
		return err
	}

	mozaikJSON, err := mozaik.ToJSON()
	if err != nil {
		return err
	}

	err = os.WriteFile(path.Join(mozaikPath, "mozaik.json"), mozaikJSON, 0644)
	if err != nil {
		return err
	}

	if mozaik.OriginalImage != nil {
		originalImage, err := os.Create(path.Join(mozaikPath, "original.png"))
		if err != nil {
			return err
		}
		defer originalImage.Close()

		err = png.Encode(originalImage, mozaik.OriginalImage)
		if err != nil {
			return err
		}
	}

	return nil
}
