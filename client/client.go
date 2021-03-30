package client

import (
	"errors"
	"fmt"
	"github.com/karalabe/hid"
	"go-streamdeck-interface/button_image"
	"go-streamdeck-interface/devices"
)


const vendorID = 0x0fd9

type Client struct {
	fd *hid.Device
	deviceInfo devices.DeviceInfo
	buttonDown []bool
}

// Open a Streamdeck device, the most common entry point
func Open() (*Client, error) {
	usbDevices := hid.Enumerate(vendorID, 0)
	if len(usbDevices) == 0 {
		return nil, errors.New("no elgato devices found")
	}

	retval := &Client{}

	found := false
	var err error
	for _, dev := range usbDevices {
		if di, ok := devices.AllDeviceTypes[dev.ProductID]; ok {
			found = true
			retval.fd, err = dev.Open()
			if err != nil {
				return nil, err
			}
			retval.deviceInfo = di
			break
		}
	}
	if !found {
		return nil, errors.New("found an Elgato device, but not one for which there is a definition; the device is not supported")
	}
	err = retval.ResetComms()
	if err != nil {
		return nil, err
	}
	retval.buttonDown = make([]bool, retval.deviceInfo.GetDeviceInfoData().GetTotalNumButtons())
	return retval, nil
}

// Close the device
func (c *Client) Close() {
	c.fd.Close()
}

// ResetComms will reset the comms protocol to the StreamDeck; useful if things have gotten de-synced, but it will also reboot the StreamDeck
func (c *Client) ResetComms() error {
	payload := c.deviceInfo.GetDeviceInfoData().ResetPacket
	_, err := c.fd.SendFeatureReport(payload)
	if err != nil {
		return err
	}
	return nil
}

func (c* Client) WriteButtonImage(btnIndex uint, i* button_image.ButtonImage) error {
	if btnIndex >= c.deviceInfo.GetDeviceInfoData().GetTotalNumButtons() {
		return errors.New(fmt.Sprintf("Invalid key index: %d", btnIndex))
	}

	bytesRemaining := len(i.Bytes)
	var pageNumber int = 0

	for bytesRemaining > 0 {
		header := c.deviceInfo.GetImageHeader(uint(bytesRemaining), uint(btnIndex), uint(pageNumber))
		imageReportLength := int(c.deviceInfo.GetDeviceInfoData().ImageReportPayloadLength)
		imageReportHeaderLength := len(header)
		imageReportPayloadLength := imageReportLength - imageReportHeaderLength

		thisLength := 0
		if imageReportPayloadLength < bytesRemaining {
			thisLength = imageReportPayloadLength
		} else {
			thisLength = bytesRemaining
		}

		bytesSent := pageNumber * imageReportPayloadLength

		payload := append(header, i.Bytes[bytesSent:(bytesSent+thisLength)]...)
		padding := make([]byte, imageReportLength-len(payload))

		thingToSend := append(payload, padding...)
		_, err := c.fd.Write(thingToSend)
		if err != nil {
			return err
		}

		bytesRemaining = bytesRemaining - thisLength
		pageNumber = pageNumber + 1
	}
	return nil
}

// SetBrightness sets the button brightness
// pct is an integer between 0-100
func (c* Client) SetBrightness(pct int) error {
	if pct < 0 {
		pct = 0
	}
	if pct > 100 {
		pct = 100
	}

	preamble := c.deviceInfo.GetDeviceInfoData().BrightnessPacket
	payload := append(preamble, byte(pct))
	_, err := c.fd.SendFeatureReport(payload)
	if err != nil {
		return err
	}
	return nil
}

func (c* Client) ReadNextButtonEvents() (downButtons []uint, upButtons []uint, err error) {
	for {
		did := c.deviceInfo.GetDeviceInfoData()
		data := make([]byte, did.GetTotalNumButtons()+did.ButtonReadOffset)
		_, err := c.fd.Read(data)
		if err != nil {
			return
		}
		for i := uint(0); i <did.GetTotalNumButtons(); i++ {
			if data[did.ButtonReadOffset+i] == 1 {
				if !c.buttonDown[i] {
					downButtons = append(downButtons, i)
				}
				c.buttonDown[i]= true
			} else {
				if c.buttonDown[i] {
					upButtons = append(upButtons, i)
				}
				c.buttonDown[i] = false
			}
		}
	}
}

type ButtonEvent struct {
	Id uint
	Down bool
}

func (c *Client) ButtonEventsChannel() (chan ButtonEvent, chan error) {
	eventChannel := make(chan ButtonEvent)
	errorChannel := make(chan error)

	go func() {
		for {
			upEvents, downEvents, err := c.ReadNextButtonEvents()
			if err != nil {
				errorChannel <- err
				return
			}
			for _, event := range downEvents {
				eventChannel <- ButtonEvent {
					event,
					true,
				}
			}
			for _, event := range upEvents {
				eventChannel <- ButtonEvent {
					event,
					false,
				}
			}
		}
	}()

	return eventChannel, errorChannel
}
