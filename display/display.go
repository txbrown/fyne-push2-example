package display

import (
	"context"
	"errors"
	"log"

	"github.com/google/gousb"
	"github.com/google/gousb/usbid"

	"flag"
	"fmt"
)

var abletonVendorID gousb.ID = 0x2982

var pushProductID gousb.ID = 0x1967

var frameHeader = []byte{0xff, 0xcc, 0xaa, 0x88,
	0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00}

var (
	debug = flag.Int("debug", 0, "libusb debug level (0..3)")
)

// AbletonPush2Display - Ableton Push 2 Display Interface
type AbletonPush2Display struct {
	pixels []byte
	device *gousb.Device
	ctx    *gousb.Context
	intf   *gousb.Interface
}

// NewAbletonPush2Display - returns a new instance
func NewAbletonPush2Display() AbletonPush2Display {
	return AbletonPush2Display{}
}

// Close - Releases handle to device and usb context
func (d *AbletonPush2Display) Close() error {
	if d.device == nil || d.ctx == nil {
		return nil
	}

	err := d.device.Close()

	if err != nil {
		panic(err)
	}

	err = d.ctx.Close()

	if err != nil {
		panic(err)
	}

	return nil
}

// Open - Opens the Ableton Push 2 usb device
func (d *AbletonPush2Display) Open() error {
	// Only one context should be needed for an application.  It should always be closed.
	ctx := gousb.NewContext()
	defer ctx.Close()

	// Debugging can be turned on; this shows some of the inner workings of the libusb package.
	ctx.Debug(*debug)

	device, err := ctx.OpenDeviceWithVIDPID(abletonVendorID, pushProductID)

	// OpenDevices can occasionally fail, so be sure to check its return value.
	if err != nil {
		log.Fatalf("list: %s", err)
		return err
	}

	if device == nil {
		return errors.New("Unable to open device")
	}

	log.Println("Device opened")

	fmt.Printf("%03d.%03d %s:%s %s\n", device.Desc.Bus, device.Desc.Address, device.Desc.Vendor, device.Desc.Product, usbid.Describe(device.Desc))

	fmt.Printf("  Protocol: %s\n", usbid.Classify(device.Desc))

	interF, _, err := device.DefaultInterface()

	d.intf = interF

	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

type contextReader interface {
	ReadContext(context.Context, []byte) (int, error)
}

// WritePixels - sends pixels to Ableton Push 2 display
func (d *AbletonPush2Display) WritePixels(pixels []uint8) error {

	outEp, _ := d.intf.OutEndpoint(1)

	_, err := outEp.Write(frameHeader)

	if err != nil {
		return err
	}

	_, err = outEp.Write(pixels)

	if err != nil {
		return err
	}

	return nil
}
