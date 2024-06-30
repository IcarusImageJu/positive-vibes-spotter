package main

import (
	"fmt"
	"os"
	"time"

	"positive-vibes-spotter/caption"
	"positive-vibes-spotter/capture"
	logger "positive-vibes-spotter/log"
	"positive-vibes-spotter/render"
	"positive-vibes-spotter/utils"

	"github.com/go-resty/resty/v2"
)

var (
	logFile   string
	apiKey    string
	model     = "gpt-4o"
	imagePath = "photo.jpg"
)

func init() {
	defer func() {
		if r := recover(); r != nil {
			logger.Error("Une erreur inattendue s'est produite: ", r)
		}
		logger.Info("Fin de l'exécution du script à ", time.Now())
	}()

	// Initialize logging
	logDir := "/spot/logs"
	logger.Info("Création du répertoire de logs: ", logDir)
	err := os.MkdirAll(logDir, os.ModePerm)
	if err != nil {
		logger.Fatal("Erreur lors de la création du répertoire de logs: ", err)
	}

	// Vérifier les permissions du répertoire
	logger.Info("Vérification des permissions du répertoire de logs")
	fileInfo, err := os.Stat(logDir)
	if err != nil {
		logger.Fatal("Erreur lors de la vérification des permissions du répertoire: ", err)
	}
	logger.Printf("Permissions du répertoire: %s", fileInfo.Mode())

	logFile = fmt.Sprintf("%s/script_%s.log", logDir, time.Now().Format("2006-01-02_15-04-05"))
	logger.Info("Fichier de log: ", logFile)
	err = logger.Setup(logFile)
	if err != nil {
		logger.Fatal("Erreur lors de la configuration du logging: ", err)
	} else {
		logger.Info("Logging configuré avec succès.")
	}
	logger.Info("Début de l'exécution du script à ", time.Now())
}

func main() {
	defer func() {
		if r := recover(); r != nil {
			logger.Error("Une erreur inattendue s'est produite: ", r)
		}
		logger.Info("Fin de l'exécution du script à ", time.Now())
	}()

	// Load environment variables from .env file
	err := utils.LoadEnv(".env")
	if err != nil {
		logger.Fatal("Erreur lors du chargement des variables d'environnement: ", err)
	}

	apiKey = os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		logger.Error("Erreur : La clé API n'est pas définie.")
		os.Exit(1)
	}

	// Check and install necessary commands
	capture.CheckInstall("convert", "imagemagick")
	capture.CheckInstall("jq", "jq")
	capture.CheckInstall("fim", "fim")
	capture.CheckInstall("libcamera-jpeg", "libcamera-apps")
	capture.CheckAndInstallFonts()

	// Take a photo
	capture.TakePhoto(imagePath)

	// Encode image to base64
	imageBase64 := capture.EncodeImageToBase64(imagePath)

	// Create content
	content := caption.CreateContent()

	// Create JSON payload
	requestPayload := caption.CreatePayload(content, imageBase64, model)

	// Send request to OpenAI API
	logger.Info("Envoi de la requête à l'API OpenAI")
	client := resty.New()
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", "Bearer "+apiKey).
		SetBody(requestPayload).
		Post("https://api.openai.com/v1/chat/completions")
	if err != nil {
		logger.Fatal("Erreur lors de l'envoi de la requête: ", err)
	}

	// Log raw response
	logger.Info("Réponse brute: ", string(resp.Body()))

	// Extract caption from response
	captionText := caption.ExtractCaption(resp.Body())

	// Create image with caption
	render.CreateImageWithCaption(captionText, "output.png")

	// Display image
	render.DisplayImage("output.png")
}