#!/bin/bash
# Déploiement automatique du programme Go sur le Raspberry Pi

# Variables
source .env

# Vérification de la compilation du programme Go en local
echo "Vérification de la compilation du programme Go en local..."
go build -o /dev/null main.go

if [ $? -ne 0 ]; then
    echo "Échec de la compilation locale. Déploiement annulé."
    exit 1
fi

# Synchronisation des fichiers
echo "Synchronisation des fichiers..."
rsync -avz --exclude='.git' $LOCAL_PATH $PI_USER@$PI_HOST:$PI_PATH

# Arrêt du programme en cours sur le Raspberry Pi
echo "Arrêt des processus 'go run main.go' en cours sur le Raspberry Pi..."
ssh $PI_USER@$PI_HOST "pkill -f 'go run main.go'"

# Pause pour assurer que les processus sont arrêtés
sleep 2

# Vérification des processus arrêtés
if ssh $PI_USER@$PI_HOST "pgrep -f 'go run main.go' > /dev/null"; then
    echo "Les processus 'go run main.go' sont toujours en cours. Tentative d'arrêt forcé..."
    ssh $PI_USER@$PI_HOST "pkill -9 -f 'go run main.go'"
    sleep 2
    if ssh $PI_USER@$PI_HOST "pgrep -f 'go run main.go' > /dev/null"; then
        echo "Impossible d'arrêter les processus 'go run main.go'. Déploiement annulé."
        exit 1
    else
        echo "Processus arrêtés avec succès après l'arrêt forcé."
    fi
else
    echo "Tous les processus 'go run main.go' arrêtés avec succès."
fi

# Exécution du script sur le Raspberry Pi
echo "Exécution du script sur le Raspberry Pi..."
ssh $PI_USER@$PI_HOST "cd $PI_PATH && nohup go run main.go > /dev/null 2>&1 &"

# Vérification que le nouveau processus est lancé
sleep 2
if ssh $PI_USER@$PI_HOST "pgrep -f 'go run main.go' > /dev/null"; then
    echo "Nouveau processus lancé avec succès"
else
    echo "Échec du lancement du nouveau processus"
    exit 1
fi