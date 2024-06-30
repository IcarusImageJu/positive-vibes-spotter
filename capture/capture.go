package capture

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os/exec"
	logger "positive-vibes-spotter/log"
)

// CheckInstall verifies if a command is installed, and installs the corresponding package if not.
func CheckInstall(command string, pkg string) {
	_, err := exec.LookPath(command)
	if err != nil {
		logger.Info(fmt.Sprintf("%s n'est pas installé. Installation...", command))

		cmd := exec.Command("sudo", "apt-get", "install", pkg, "-y")

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
			logger.Fatal(fmt.Sprintf("Erreur lors de l'installation de %s: %v", pkg, err))
		}
	} else {
		logger.Info(fmt.Sprintf("%s est déjà installé.", command))
	}
}

// CheckAndInstallFonts ensures Arial fonts are installed, and installs them if not.
func CheckAndInstallFonts() {
	arialInstalled := exec.Command("fc-list")
	arialInstalledOut, err := arialInstalled.Output()
	if err != nil {
		logger.Fatal(fmt.Sprintf("Erreur lors de la vérification des polices: %v", err))
	}
	if !bytes.Contains(bytes.ToLower(arialInstalledOut), []byte("arial")) {
		logger.Info("Les polices Arial ne sont pas installées. Installation...")

		cmd := exec.Command("sudo", "apt-get", "install", "ttf-mscorefonts-installer", "-y")

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
			logger.Fatal(fmt.Sprintf("Erreur lors de l'installation des polices Arial: %v", err))
		}

		cmd = exec.Command("sudo", "fc-cache", "-f", "-v")
		cmd.Stdout = stdout
		cmd.Stderr = stderr
		if err := cmd.Run(); err != nil {
			logger.Fatal(fmt.Sprintf("Erreur lors de la mise en cache des polices: %v", err))
		}
	} else {
		logger.Info("Les polices Arial sont déjà installées.")
	}
}

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
	imageBytes, err := ioutil.ReadFile(imagePath)
	if err != nil {
		logger.Fatal(fmt.Sprintf("Erreur lors de la lecture de l'image: %v", err))
	}
	logger.Info("Image encodée en base64 avec succès.")
	return base64.StdEncoding.EncodeToString(imageBytes)
}