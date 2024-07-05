#!/bin/bash
# Déploiement automatique du programme Go sur le Raspberry Pi

# Variables
source .env

# Synchronisation des fichiers
echo "Synchronisation des fichiers..."
rsync -avz --exclude='.git' $LOCAL_PATH $PI_USER@$PI_HOST:$PI_PATH

# Arrêt du programme compilé en cours sur le Raspberry Pi
echo "Arrêt des processus 'spot' en cours sur le Raspberry Pi..."
ssh $PI_USER@$PI_HOST "pkill -f 'spot'"

# Pause pour assurer que les processus sont arrêtés
sleep 2

# Compilation du programme sur le Raspberry Pi
echo "Compilation du programme sur le Raspberry Pi..."
ssh $PI_USER@$PI_HOST "cd $PI_PATH && go build -o spot main.go"

# Vérification de la compilation
if ssh $PI_USER@$PI_HOST "[ ! -f $PI_PATH/spot ]"; then
    echo "Échec de la compilation sur le Raspberry Pi. Déploiement annulé."
    exit 1
fi

# Exécution du programme compilé sur le Raspberry Pi
echo "Exécution du programme compilé sur le Raspberry Pi..."
ssh $PI_USER@$PI_HOST "cd $PI_PATH && ./spot" &

# Vérification que le nouveau processus est lancé
sleep 2
if ssh $PI_USER@$PI_HOST "pgrep -f './spot' > /dev/null"; then
    echo "Nouveau processus lancé avec succès"
else
    echo "Échec du lancement du nouveau processus"
    exit 1
fi