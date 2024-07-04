#!/bin/bash

# Variables
source .env

# Fonction de synchronisation et exécution
echo "Synchronisation des fichiers..."
rsync -av --exclude='.git' $LOCAL_PATH $PI_USER@$PI_HOST:$PI_PATH

echo "Exécution du script sur le Raspberry Pi..."
ssh $PI_USER@$PI_HOST "cd $PI_PATH && go run main.go" &