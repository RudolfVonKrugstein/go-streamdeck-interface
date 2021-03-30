package devices

import "github.com/disintegration/gift"

var (
	xlName                     string
	xlButtonWidth              uint
	xlButtonHeight             uint
	xlImageReportPayloadLength uint
)

type XLDeviceInfo struct {
}


func (xldi* XLDeviceInfo) GetDeviceInfoData() *DeviceInfoData {
	return &DeviceInfoData{
		"Streamdeck XL",
		96,
		96,
		1024,
		0,
		0,
		0x6c,
		resetPacket32(),
		4,
		8,
		brightnessPacket32(),
		4,
		JPEG,
	}
}

// GetImageHeader returns the USB comms header for a button image for the XL
func (xldi* XLDeviceInfo)  GetImageHeader(bytesRemaining uint, btnIndex uint, pageNumber uint) []byte {
	thisLength := uint(0)
	di := xldi.GetDeviceInfoData()
	if di.ImageReportPayloadLength < bytesRemaining {
		thisLength = di.ImageReportPayloadLength
	} else {
		thisLength = bytesRemaining
	}
	header := []byte{'\x02', '\x07', byte(btnIndex)}
	if thisLength == bytesRemaining {
		header = append(header, '\x01')
	} else {
		header = append(header, '\x00')
	}

	header = append(header, byte(thisLength&0xff))
	header = append(header, byte(thisLength>>8))

	header = append(header, byte(pageNumber&0xff))
	header = append(header, byte(pageNumber>>8))

	return header
}

func (xldi* XLDeviceInfo) GiftTransform() *gift.GIFT {
	return gift.New(
		gift.Resize(
			int(xldi.GetDeviceInfoData().ButtonWidth),
			int(xldi.GetDeviceInfoData().ButtonWidth), gift.LanczosResampling),
		gift.Rotate180(),
	)
}
