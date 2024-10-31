package render

import (
	"fmt"
	"os"
	"os/exec"
	logger "positive-vibes-spotter/log"
	"positive-vibes-spotter/utils"
	"time"
)

var (
	tickerTime  = 5 * time.Second
	outputPath1 = "spot1.png"
	outputPath2 = "spot2.png"
	scaleFactor = 0.625
	squareSize = 2
)

type CheckboardType string

const (
	CheckboardEven CheckboardType = "even"
	CheckboardOdd  CheckboardType = "odd"
)

// CreateImageWithCaptionAndMask crée une image avec la légende spécifiée et applique un masque
// qui rend les pixels pairs ou impairs noirs, en fonction du type de masque.
func CreateImageWithCaptionAndMask(caption string, outputPath string, mask CheckboardType) {
	logger.Info("Création de l’image avec la légende et le masque")

	stdout, err := logger.Writer()
	if err != nil {
		logger.Fatal(fmt.Sprintf("Erreur lors de l'ouverture du fichier de log pour stdout: %v", err))
	}
	defer stdout.Close()

	stderr, err := logger.Writer()
	if err != nil {
		logger.Fatal(fmt.Sprintf("Erreur lors de l'ouverture du fichier de log pour stderr: %v", err))
	}
	defer stderr.Close()

	// Calculer les dimensions en fonction du scaleFactor
	baseWidth := int(1280 * scaleFactor)
	baseHeight := int(720 * scaleFactor)
	border := int(100 * scaleFactor)
	extentOffset := int(50 * scaleFactor)
	pointSize := int(72 * scaleFactor)
	captionSizeWidth := int(1200 * scaleFactor)
	captionSizeHeight := int(600 * scaleFactor)

	// Créer l'image initiale avec la légende
	cmd := exec.Command("convert",
		"-background", "black",
		"-fill", "#888888",
		"-font", "Arial",
		"-pointsize", fmt.Sprintf("%d", pointSize),
		"-gravity", "southwest",
		"-extent", fmt.Sprintf("%dx%d", baseWidth, baseHeight),
		"-size", fmt.Sprintf("%dx%d", captionSizeWidth, captionSizeHeight),
		"caption:"+caption,
		"-bordercolor", "black",
		"-border", fmt.Sprintf("%dx%d", border, border),
		"-gravity", "southwest",
		"-extent", fmt.Sprintf("%dx%d+%d+%d", baseWidth, baseHeight, extentOffset, extentOffset),
		"original_image.png")
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	err = cmd.Run()
	if err != nil {
		logger.Fatal(err)
	}

	// Créer un masque pour les pixels pairs ou impairs avec des carrés de taille squareSize
	maskExpression := fmt.Sprintf("mod(int(floor(i/%d)+floor(j/%d)),2)", squareSize, squareSize)
	if mask == CheckboardOdd {
			maskExpression = fmt.Sprintf("1-mod(int(floor(i/%d)+floor(j/%d)),2)", squareSize, squareSize)
	}


	// Créer le masque avec ImageMagick
	maskCmd := exec.Command("convert",
		"-size", fmt.Sprintf("%dx%d", baseWidth, baseHeight),
		"xc:black",
		"-fill", "white",
		"-fx", maskExpression,
		"mask.png")
	maskCmd.Stdout = stdout
	maskCmd.Stderr = stderr
	err = maskCmd.Run()
	if err != nil {
		logger.Fatal(err)
	}

	// Appliquer le masque à l'image originale
	applyMaskCmd := exec.Command("convert",
		"original_image.png",
		"mask.png",
		"-compose", "Multiply",
		"-composite",
		outputPath)
	applyMaskCmd.Stdout = stdout
	applyMaskCmd.Stderr = stderr
	err = applyMaskCmd.Run()
	if err != nil {
		logger.Fatal(err)
	}

	// Supprimer les fichiers temporaires
	os.Remove("original_image.png")
	os.Remove("mask.png")
}

var previousCmd *exec.Cmd

func display(outputPath string) {
	logger.Info("Affichage de l’image")

	// Terminer le processus précédent si nécessaire
	if previousCmd != nil && previousCmd.Process != nil {
		err := previousCmd.Process.Kill()
		if err != nil {
			logger.Error(fmt.Sprintf("Erreur lors de la fermeture de l'instance précédente de fim: %v", err))
		}
	}

	cmd := exec.Command("fim", "-a", "--quiet", outputPath)
	err := cmd.Start()
	if err != nil {
		logger.Fatal(err)
	}

	previousCmd = cmd
}

func cleanup() {
	if previousCmd != nil && previousCmd.Process != nil {
		err := previousCmd.Process.Kill()
		if err != nil {
			logger.Error(fmt.Sprintf("Erreur lors de la fermeture de fim: %v", err))
		}
	}
}

func alternateImage(image1, image2 string) {
	ticker := time.NewTicker(tickerTime)
	defer func() {
		ticker.Stop()
		cleanup()
	}()

	displayFirst := true
	for range ticker.C {
		if displayFirst {
			display(image1)
		} else {
			display(image2)
		}
		displayFirst = !displayFirst
	}
}

func Render(caption string) {
	utils.CheckInstall("fim", "fim")
	utils.CheckInstall("convert", "imagemagick")
	utils.CheckAndInstallFonts()
	CreateImageWithCaptionAndMask(caption, outputPath1, CheckboardEven)
	CreateImageWithCaptionAndMask(caption, outputPath2, CheckboardOdd)

	alternateImage(outputPath1, outputPath2)
}
