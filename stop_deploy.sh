#!/bin/bash

# Vérifier les processus en cours avec pgrep
deploy_pids=$(pgrep -f "ssh $PI_USER@$PI_HOST")

if [ -z "$deploy_pids" ]; then
    echo "Aucun processus de déploiement en cours trouvé."
else
    # Tuer les processus identifiés
    echo "Arrêt des processus de déploiement en cours..."
    echo "$deploy_pids" | xargs kill
    echo "Processus de déploiement arrêtés."
fi