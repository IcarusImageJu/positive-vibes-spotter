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
)

var (
	logFile   string
	apiKey    = utils.GetEnv("OPENAI_API_KEY")
	model     = "gpt-4o"
	imagePath = "photo.jpg"
	outputPath = "output.png"
	tickerTime	= 1 * time.Hour
)

func init() {
	// Setup defer to handle any panics and log the script end time
	defer func() {
		if r := recover(); r != nil {
			logger.Error("Une erreur inattendue s'est produite: ", r)
		}
		logger.Info("Fin de l'exécution du script à ", time.Now())
	}()

	// Initialize logging
	logDir := "/spot/logs"
	logger.Info("Création du répertoire de logs: ", logDir)
	if err := os.MkdirAll(logDir, os.ModePerm); err != nil {
		logger.Fatal("Erreur lors de la création du répertoire de logs: ", err)
	}

	logFile = fmt.Sprintf("%s/script_%s.log", logDir, time.Now().Format("2006-01-02_15-04-05"))
	logger.Info("Fichier de log: ", logFile)
	if err := logger.Setup(logFile); err != nil {
		logger.Fatal("Erreur lors de la configuration du logging: ", err)
	} else {
		logger.Info("Logging configuré avec succès.")
	}

	logger.Info("Début de l'exécution du script à ", time.Now())
}

func main() {
	ticker := time.NewTicker(tickerTime)
	defer ticker.Stop()

	// Run the task immediately before starting the ticker
	runTask()

	for range ticker.C {
		runTask()
	}
}

func runTask() {
	// Setup defer to handle any panics and log the script end time
	defer func() {
		if r := recover(); r != nil {
			logger.Error("Une erreur inattendue s'est produite: ", r)
		}
		logger.Info("Fin de l'exécution du script à ", time.Now())
	}()

	// Take a photo
	picture := capture.Picture(imagePath)

	// Get caption
	caption := caption.Caption(picture, model, apiKey)

	// Render
	render.Render(caption, picture, outputPath)
}