package render

import (
	"image"
	"image/color"
	"image/draw"
	"log"
	"math"
	"math/rand"
	"os"
	"time"
)

const (
	width      = 800
	height     = 480
	tickerTime = 1000 / 20 * time.Millisecond
)

func rgbaToRGB565(c color.RGBA) uint16 {
	r := uint16(c.R >> 3)
	g := uint16(c.G >> 2)
	b := uint16(c.B >> 3)
	return (r << 11) | (g << 5) | b
}

func writeToFramebuffer(img *image.RGBA, file *os.File) error {
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

	_, err := file.WriteAt(buffer, 0)
	return err
}

func drawSine(amplitudeMultiplier, frequencyMultiplier, phase1, phase2, phase3, dynamicHeight, offsetY float64, x int) int {
	// Calcul de y1 avec les trois composantes sinusoidales
	y1 := amplitudeMultiplier * math.Sin(frequencyMultiplier*float64(x)/float64(width)*4*math.Pi+phase1) * 20 * (0.5 + 0.5*math.Log1p(float64(height-dynamicHeight)/height))
	y2 := amplitudeMultiplier * math.Sin(frequencyMultiplier*float64(x)/float64(width)*8*math.Pi+phase2/2) * 10 * (0.5 + 0.5*math.Log1p(float64(height-dynamicHeight)/height))
	y3 := amplitudeMultiplier * math.Sin(frequencyMultiplier*float64(x)/float64(width)*2*math.Pi+phase3/2) * 5 * (0.5 + 0.5*math.Log1p(float64(height-dynamicHeight)/height))
	y := int(y1 + y2 + y3 + float64(dynamicHeight) + float64(offsetY))
	return y
}

func drawWave(img *image.RGBA, phase1, phase2, phase3, amplitudeMultiplier, frequencyMultiplier, dynamicHeight, offsetY float64) {
	white := color.RGBA{255, 255, 255, 255}

	for x := 0; x < width; x++ {
		// Calcul de y1 avec les trois composantes sinusoidales
		y1 := drawSine(amplitudeMultiplier, frequencyMultiplier, phase1, phase2, phase3, dynamicHeight, offsetY, x)
		y2 := drawSine(amplitudeMultiplier, frequencyMultiplier, phase2, phase3, phase1, dynamicHeight, offsetY, x)
		y3 := drawSine(amplitudeMultiplier, frequencyMultiplier, phase3, phase1, phase2, dynamicHeight, offsetY, x)

		// Calcul des autres courbes
		// y2 := int(amplitudeMultiplier * math.Sin(frequencyMultiplier*float64(x)/float64(width)*4*math.Pi+phase2+math.Pi/2)*20 + dynamicHeight + offsetY)
		// y3 := int(amplitudeMultiplier * math.Sin(frequencyMultiplier*float64(x)/float64(width)*4*math.Pi+phase3+math.Pi)*20 + dynamicHeight + offsetY)

		// Dessin de y1
		if y1 >= 0 && y1 < height {
			img.Set(x, y1, white)
			img.Set(x+1, y1+1, white)
			img.Set(x+2, y1+2, white)
			img.Set(x+3, y1+2, white)
			img.Set(x+4, y1+2, white)
			img.Set(x+5, y1+2, white)
			img.Set(x+6, y1+2, white)
			img.Set(x+2, y1+3, white)
			img.Set(x+3, y1+3, white)
			img.Set(x+4, y1+3, white)
			img.Set(x+5, y1+3, white)
			img.Set(x+6, y1+3, white)
		}

		// Dessin de y2 et y3 seulement s'ils ne sont pas en dessous de y1
		if y2-20 >= 0 && y2-20 < height && y2-20 <= y1 {
			img.Set(x, y2-20, white)
		}
		if y3-25 >= 0 && y3-25 < height && y3-25 <= y1 {
			img.Set(x, y3-25, white)
		}

		y := 0.0
		for i := 0; i < 5; i++ {
			frequency := rand.Float64() * 3 * math.Pi
			amplitude := rand.Float64() * 200
			offset := rand.Float64() * 2 * math.Pi
			y += amplitude * math.Sin(frequency*float64(x)/float64(width)+phase1+offset)
		}
		y = y/5 + dynamicHeight - 50

		// Dessin du bruit seulement s'il n'est pas en dessous de y1
		if int(y)+int(offsetY) >= 0 && int(y)+int(offsetY) < height && int(y)+int(offsetY) <= y1 {
			img.Set(x, int(y)+int(offsetY), white)
		}

	}
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

	// Variables pour gérer les vagues
	phase1_1 := 0.0
	phase1_2 := 0.0
	phase1_3 := 0.0
	phase2_1 := 0.0
	phase2_2 := 0.0
	phase2_3 := 0.0
	dynamicHeight1 := -60.0
	dynamicHeight2 := -60.0
	speed1 := 1.0
	speed2 := 1.0
	rand.Seed(time.Now().UnixNano())

	for range ticker.C {
		draw.Draw(img, img.Bounds(), &image.Uniform{color.Black}, image.Point{}, draw.Src)

		// Dessiner la première série de vagues
		drawWave(img, phase1_1, phase1_2, phase1_3, 1.3, 1.0, dynamicHeight1, 0)

		// Dessiner la deuxième série de vagues décalées de 50 % de la hauteur
		drawWave(img, phase2_1, phase2_2, phase2_3, 2.0, 1.1, dynamicHeight2, height/2)

		phase1_1 += 0.002
		phase1_2 += 0.12
		phase1_3 += 0.05
		phase2_1 += 0.003
		phase2_2 += 0.09
		phase2_3 += 0.1
		dynamicHeight1 += speed1
		dynamicHeight2 += speed2
		speed1 = 0.7 + 0.3*math.Log1p(float64(height-dynamicHeight1)/height)
		speed2 = 0.7 + 0.3*math.Log1p(float64(height-dynamicHeight2)/height)

		if dynamicHeight1 >= float64(height+100) {
			dynamicHeight1 = -100
			speed1 = 1.0 // Réinitialise la vitesse pour la prochaine vague
			phase1_1 = 0.0
			phase1_2 = 0.0
			phase1_3 = 0.0
		}
		if dynamicHeight2 >= float64(height+100 - height/2) {
			dynamicHeight2 = -100 - height/2
			speed2 = 1.0 // Réinitialise la vitesse pour la prochaine vague
			phase2_1 = 0.0
			phase2_2 = 0.0
			phase2_3 = 0.0
		}

		if err := writeToFramebuffer(img, file); err != nil {
			log.Fatalf("Failed to write to framebuffer: %v", err)
		}
	}

	return nil
}