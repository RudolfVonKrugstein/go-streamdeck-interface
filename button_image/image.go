package button_image

import (
	"bytes"
	"errors"
	"github.com/RudolfVonKrugstein/go-streamdeck-interface/devices"
	"golang.org/x/image/bmp"
	"image"
	"image/color"
	"image/jpeg"
)

type ButtonImage struct {
	Bytes []byte
}

func NewButtonImage(di devices.DeviceInfo, img image.Image) (*ButtonImage, error) {
	// First device specific rotation
	g := di.GiftTransform()
	res := image.NewRGBA(g.Bounds(img.Bounds()))
	g.Draw(res, img)

	// Now encode in the specific format
	var b bytes.Buffer
	switch di.GetDeviceInfoData().ImageFormat {
	case devices.JPEG:
		jpeg.Encode(&b, img, &jpeg.Options{Quality: 100})
	case devices.BMP:
		bmp.Encode(&b, img)
	default:
		return nil, errors.New("unknown button image format")
	}
	return &ButtonImage{b.Bytes()}, nil
}

func NewSolidColorButtonImage(di devices.DeviceInfo, color color.Color) (*ButtonImage, error) {
	return NewButtonImage(di, image.NewUniform(color))
}
