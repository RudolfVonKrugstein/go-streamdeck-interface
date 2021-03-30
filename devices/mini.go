package devices

import "github.com/disintegration/gift"

type MiniDeviceInfo struct {
}

func (mdi* MiniDeviceInfo) GetDeviceInfoData() *DeviceInfoData {
	return &DeviceInfoData{
		"Streamdeck Mini",
		80,
		80,
		1024,
		0,
		0,
		0x63,
		resetPacket17(),
		2,
		3,
		brightnessPacket17(),
		1,
		BMP,
	}
}


// GetImageHeaderMini returns the USB comms header for a button image for the Mini
func (mdi* MiniDeviceInfo) GetImageHeader(bytesRemaining uint, btnIndex uint, pageNumber uint) []byte {
	var thisLength uint
	di := mdi.GetDeviceInfoData()
	if di.ImageReportPayloadLength < bytesRemaining {
		thisLength = di.ImageReportPayloadLength
	} else {
		thisLength = bytesRemaining
	}
	header := []byte{
		'\x02',
		'\x01',
		byte(pageNumber),
		0,
		get_header_element(thisLength, bytesRemaining),
		byte(btnIndex + 1),
		'\x00',
		'\x00',
		'\x00',
		'\x00',
		'\x00',
		'\x00',
		'\x00',
		'\x00',
		'\x00',
		'\x00',
	}
	return header
}

func (mdi* MiniDeviceInfo) GiftTransform() *gift.GIFT {
	return gift.New(
		gift.Resize(
			int(mdi.GetDeviceInfoData().ButtonWidth),
			int(mdi.GetDeviceInfoData().ButtonWidth), gift.LanczosResampling),
		gift.Rotate90(),
		gift.FlipVertical(),
	)
}

func get_header_element(thisLength, bytesRemaining uint) byte {
	if thisLength == bytesRemaining {
		return '\x01'
	} else {
		return '\x00'
	}
}
