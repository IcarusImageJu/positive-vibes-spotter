package render

import (
	"image"
	"image/color"
	"image/draw"
	"log"
	"os"
	"time"
)

const (
	width      = 800
	height     = 480
	tickerTime = 1 * time.Second
)

func rgbaToRGB565(c color.RGBA) uint16 {
	r := uint16(c.R >> 3)
	g := uint16(c.G >> 2)
	b := uint16(c.B >> 3)
	return (r << 11) | (g << 5) | b
}

func writeToFramebuffer(img *image.RGBA, file *os.File) error {
	// Create a buffer for the framebuffer data
	buffer := make([]byte, width*height*2)

	for iy := 0; iy < height; iy++ {
		for ix := 0; ix < width; ix++ {
			offset := (iy*img.Stride + ix*4)
			rgba := color.RGBA{
				R: img.Pix[offset+0],
				G: img.Pix[offset+1],
				B: img.Pix[offset+2],
				A: img.Pix[offset+3],
			}
			rgb565 := rgbaToRGB565(rgba)
			buffer[(iy*width+ix)*2] = byte(rgb565 & 0xff)
			buffer[(iy*width+ix)*2+1] = byte((rgb565 >> 8) & 0xff)
		}
	}

	// Write the entire buffer to the framebuffer in one go
	_, err := file.WriteAt(buffer, 0)
	return err
}

func Render(caption string) error {
	file, err := os.OpenFile("/dev/fb0", os.O_RDWR, 0600)
	if err != nil {
		log.Fatalf("Failed to open framebuffer device: %v", err)
	}
	defer file.Close()

	img := image.NewRGBA(image.Rect(0, 0, width, height))

	ticker := time.NewTicker(tickerTime)
	defer ticker.Stop()

	for range ticker.C {
		// Clear the image with a black background
		draw.Draw(img, img.Bounds(), &image.Uniform{color.Black}, image.Point{}, draw.Src)

		// Draw a large white square at the center
		squareSize := 300
		startX := (width - squareSize) / 2
		startY := (height - squareSize) / 2
		white := color.RGBA{255, 255, 255, 255}

		for y := startY; y < startY+squareSize; y++ {
			for x := startX; x < startX+squareSize; x++ {
				img.Set(x, y, white)
			}
		}

		if err := writeToFramebuffer(img, file); err != nil {
			log.Fatalf("Failed to write to framebuffer: %v", err)
		}
	}

	return nil
}