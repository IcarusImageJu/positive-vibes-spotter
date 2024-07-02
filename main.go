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
	tickerTime	= 1 * time.Hour
	debug     = true
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
	// Setup defer to handle any panics and log the script end time
	defer func() {
		if r := recover(); r != nil {
			logger.Error("Une erreur inattendue s'est produite: ", r)
		}
		logger.Info("Fin de l'exécution du script à ", time.Now())
	}()

	ticker := time.NewTicker(tickerTime)
	defer ticker.Stop()

	// Run the task immediately before starting the ticker
	runTask()

	for range ticker.C {
		runTask()
	}
}

func runTask() {
	logger.Info("Début de l'exécution de la tâche à ", time.Now())
	
	// Take a photo
	var picture string
	if debug {
		picture = ""
	} else {
		picture = capture.Picture(imagePath)
	}

	// Get caption
	var spotCaption string
	if debug {
		spotCaption = "Lorem ipsum dolor sit amet, consectetur adipiscing elit."
	} else {
		spotCaption = caption.Caption(picture, model, apiKey)
	}

	// Render
	render.Render(spotCaption)

	logger.Info("fin de l'exécution de la tâche à ", time.Now())
}
