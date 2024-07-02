package render

import (
	"image"
	"image/color"
	"image/draw"
	"log"
	"os"
	"strings"

	"github.com/golang/freetype"
)

const (
	width  = 800
	height = 480
)

func drawText(img *image.RGBA, text string, x, y int, col color.Color) error {
	fontBytes, err := os.ReadFile("/usr/share/fonts/truetype/dejavu/DejaVuSans-Bold.ttf")
	if err != nil {
		return err
	}
	f, err := freetype.ParseFont(fontBytes)
	if err != nil {
		return err
	}

	c := freetype.NewContext()
	c.SetDPI(72)
	c.SetFont(f)
	c.SetFontSize(32)
	c.SetClip(img.Bounds())
	c.SetDst(img)
	c.SetSrc(image.NewUniform(col))

	// Calculate the max width of text in pixels
	maxWidth := width - x*2
	lineHeight := int(c.PointToFixed(40) >> 6)

	words := strings.Fields(text)
	line := ""
	yOffset := y

	for _, word := range words {
		testLine := line + word + " "
		pt := freetype.Pt(x, yOffset+lineHeight)
		width, _ := c.DrawString(testLine, pt)
		if int(width.X>>6) > maxWidth {
			// Render the current line
			pt = freetype.Pt(x, yOffset+lineHeight)
			_, err = c.DrawString(line, pt)
			if err != nil {
				return err
			}
			line = word + " "
			yOffset += lineHeight
		} else {
			line = testLine
		}
	}
	// Render the last line
	pt := freetype.Pt(x, yOffset+lineHeight)
	_, err = c.DrawString(line, pt)
	if err != nil {
		return err
	}

	return nil
}

func rgbaToRGB565(c color.RGBA) uint16 {
	r := uint16(c.R >> 3)
	g := uint16(c.G >> 2)
	b := uint16(c.B >> 3)
	return (r << 11) | (g << 5) | b
}

func Render(caption string) {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.Draw(img, img.Bounds(), &image.Uniform{color.Black}, image.Point{}, draw.Src)

	err := drawText(img, caption, 50, 50, color.White)
	if err != nil {
		log.Fatalf("Failed to draw text: %v", err)
	}

	file, err := os.OpenFile("/dev/fb0", os.O_RDWR, 0600)
	if err != nil {
		log.Fatalf("Failed to open framebuffer device: %v", err)
	}

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			offset := (y*img.Stride + x*4)
			rgba := color.RGBA{
				R: img.Pix[offset+0],
				G: img.Pix[offset+1],
				B: img.Pix[offset+2],
				A: img.Pix[offset+3],
			}
			rgb565 := rgbaToRGB565(rgba)
			_, err := file.Write([]byte{
				byte(rgb565 & 0xff),
				byte((rgb565 >> 8) & 0xff),
			})
			if err != nil {
				log.Fatalf("Failed to write to framebuffer: %v", err)
			}
		}
	}

	err = file.Close()
	if err != nil {
		log.Fatalf("Failed to close framebuffer device: %v", err)
	}
}