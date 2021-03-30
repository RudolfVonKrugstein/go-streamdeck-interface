package devices

import "github.com/disintegration/gift"

type OriginalV2DeviceInfo struct {
}

func (ov2di* OriginalV2DeviceInfo) GetDeviceInfoData() *DeviceInfoData {
	return &DeviceInfoData{
		"Streamdeck (original v2)",
		72,
		72,
		1024,
		0,
		0,
		0x6d,
		resetPacket32(),
		3,
		5,
		brightnessPacket32(),
		4,
		JPEG,
	}
}


// GetImageHeaderOv2 returns the USB comms header for a button image for the Original v2
func (ov2di* OriginalV2DeviceInfo) GetImageHeader(bytesRemaining uint, btnIndex uint, pageNumber uint) []byte {
	thisLength := uint(0)
	di := ov2di.GetDeviceInfoData()
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

func (ov2di* OriginalV2DeviceInfo) GiftTransform() *gift.GIFT {
	return gift.New(
		gift.Resize(
			int(ov2di.GetDeviceInfoData().ButtonWidth),
			int(ov2di.GetDeviceInfoData().ButtonWidth), gift.LanczosResampling),
		gift.Rotate180(),
	)
}
