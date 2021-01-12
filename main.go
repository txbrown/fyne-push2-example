package main

import (
	"fyne-app/display"
	"image"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/widget"
)

func main() {
	runFyneApp()
}

func runFyneApp() {
	a := app.New()

	w := a.NewWindow("Hello")

	w.Resize(fyne.Size{
		Height: 160,
		Width:  960,
	})

	hello := widget.NewLabel("Hello Fyne!")

	pd := display.NewAbletonPush2Display()

	err := pd.Open()

	if err != nil {
		panic(err)
	}
	w.SetContent(widget.NewVBox(
		hello,
		// when this button is clicked we capture the screen, convert into push 2 display pixels and write to display
		widget.NewButton("Hi!", func() {
			hello.SetText("Welcome :)")
			img := w.Canvas().Capture()
			pixels := getPixels(img)

			if err != nil {
				panic(err)
			}

			pd.WritePixels(pixels)
		}),
		widget.NewSlider(0, 1000),
	))

	w.ShowAndRun()

}

// getPixels - converts an image into expected pixel format we need to send to Push 2 display
func getPixels(image image.Image) []byte {

	xOrMasks := []uint8{0xe7, 0xf3, 0xe7, 0xff}
	displayPitch := 1920 + 128
	xorOffset := 0
	pixels := make([]uint8, displayPitch*160)

	for x := 0; x < 960; x++ {
		for y := 0; y < 160; y++ {
			c := image.At(x, y)

			red, green, blue, a := c.RGBA()

			byteOffset := y*displayPitch + x*2

			pixels[byteOffset] = (uint8(red/a) >> 2) ^ xOrMasks[xorOffset]
			xorOffset = (xorOffset + 1) % 4
			pixels[byteOffset+1] =
				((uint8(blue/a) & 0xf8) | (uint8(green/a) >> 5)) ^ xOrMasks[xorOffset]
			xorOffset = (xorOffset + 1) % 4
		}
	}

	return pixels
}
