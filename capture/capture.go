package capture

import (
	"encoding/base64"
	"fmt"
	"os"
	"os/exec"
	logger "positive-vibes-spotter/log"
	"positive-vibes-spotter/utils"
)

// TakePhoto takes a photo using the libcamera-jpeg command.
func TakePhoto(imagePath string) {
	logger.Info("Prise de photo avec libcamera-jpeg")
	cmd := exec.Command("libcamera-jpeg", "-o", imagePath, "--rotation", "180")

	stdout, err := logger.Writer()
	if err != nil {
		logger.Fatal(fmt.Sprintf("Erreur lors de l'ouverture du fichier de log pour stdout: %v", err))
	}
	stderr, err := logger.Writer()
	if err != nil {
		logger.Fatal(fmt.Sprintf("Erreur lors de l'ouverture du fichier de log pour stderr: %v", err))
	}

	cmd.Stdout = stdout
	cmd.Stderr = stderr
	if err := cmd.Run(); err != nil {
		logger.Fatal(fmt.Sprintf("Erreur lors de la prise de photo: %v", err))
	} else {
		logger.Info("Photo prise avec succès.")
	}
}

// EncodeImageToBase64 reads an image from the given path and encodes it to a base64 string.
func EncodeImageToBase64(imagePath string) string {
	logger.Info("Encodage de l'image en base64")
	imageBytes, err := os.ReadFile(imagePath)
	if err != nil {
		logger.Fatal(fmt.Sprintf("Erreur lors de la lecture de l'image: %v", err))
	}
	logger.Info("Image encodée en base64 avec succès.")
	return base64.StdEncoding.EncodeToString(imageBytes)
}

func Picture(imagePath string) string {
	utils.CheckInstall("libcamera-jpeg", "libcamera-apps")
	TakePhoto(imagePath)
	return EncodeImageToBase64(imagePath)
}