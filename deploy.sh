#!/bin/bash

LOCKFILE="/tmp/deploy.lock"
PIDFILE="/tmp/deploy.pid"

# Fonction de nettoyage
cleanup() {
    rm -f "$LOCKFILE"
    rm -f "$PIDFILE"
    exit
}

# Si un autre processus est en cours, exit
if [ -e "$LOCKFILE" ]; then
    echo "Un autre processus de déploiement est en cours. Sortie."
    exit 1
fi

# Créer un fichier de verrouillage et enregistrer le PID
touch "$LOCKFILE"
echo $$ > "$PIDFILE"

# Assurer la suppression du fichier de verrouillage en cas d'arrêt ou d'interruption
trap cleanup INT TERM EXIT

# Variables
source .env

# Vérification de la compilation du programme Go en local
echo "Vérification de la compilation du programme Go en local..."
go build -o /dev/null main.go

if [ $? -ne 0 ]; then
    echo "Échec de la compilation locale. Déploiement annulé."
    cleanup
fi

# Synchronisation des fichiers
echo "Synchronisation des fichiers..."
rsync -avz --exclude='.git' $LOCAL_PATH $PI_USER@$PI_HOST:$PI_PATH

# Exécution du script sur le Raspberry Pi
echo "Exécution du script sur le Raspberry Pi..."
ssh $PI_USER@$PI_HOST "cd $PI_PATH && go run main.go" &

# Nettoyage du fichier de verrouillage
cleanup