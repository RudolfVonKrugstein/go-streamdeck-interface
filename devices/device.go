package devices

import (
	"github.com/disintegration/gift"
)

type ImageFormat int

const (
	BMP ImageFormat = 0
	JPEG ImageFormat = 1
)

var (
	AllDeviceTypes map[uint16]DeviceInfo
)

func init() {
	for _, di := range []DeviceInfo{&MiniDeviceInfo{}, &OriginalV2DeviceInfo{}, &XLDeviceInfo{}} {
		AllDeviceTypes[di.GetDeviceInfoData().UsbProductID] = di
	}
}

type DeviceInfoData struct {
	Name                     string
	ButtonWidth              uint
	ButtonHeight             uint
	ImageReportPayloadLength uint
	ImageReportHeaderLength  uint
	ImageReportLength        uint
	UsbProductID             uint16
	ResetPacket              []byte
	ButtonRows               uint
	ButtonCols               uint
	BrightnessPacket         []byte
	ButtonReadOffset         uint
	ImageFormat              ImageFormat
}

type DeviceInfo interface {
	GetDeviceInfoData() *DeviceInfoData
	GetImageHeader(bytesRemaining uint, btnIndex uint, pageNumber uint) []byte
	GiftTransform() *gift.GIFT
}

func (di* DeviceInfoData) GetTotalNumButtons() uint {
	return di.ButtonRows * di.ButtonCols
}

// resetPacket17 gives the reset packet for devices which need it to be 17 bytes long
func resetPacket17() []byte {
	pkt := make([]byte, 17)
	pkt[0] = 0x0b
	pkt[1] = 0x63
	return pkt
}

// resetPacket32 gives the reset packet for devices which need it to be 32 bytes long
func resetPacket32() []byte {
	pkt := make([]byte, 32)
	pkt[0] = 0x03
	pkt[1] = 0x02
	return pkt
}

// brightnessPacket17 gives the brightness packet for devices which need it to be 17 bytes long
func brightnessPacket17() []byte {
	pkt := make([]byte, 17)
	pkt[0] = 0x05
	pkt[1] = 0x55
	pkt[2] = 0xaa
	pkt[3] = 0xd1
	pkt[4] = 0x01
	return pkt
}

// brightnessPacket32 gives the brightness packet for devices which need it to be 32 bytes long
func brightnessPacket32() []byte {
	pkt := make([]byte, 32)
	pkt[0] = 0x03
	pkt[1] = 0x08
	return pkt
}
