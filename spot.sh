#!/bin/bash

# Répertoire pour les logs
LOG_DIR="/mood/logs"
sudo mkdir -p $LOG_DIR
LOG_FILE="$LOG_DIR/script_$(date +'%Y-%m-%d_%H-%M-%S').log"

{
    echo "Début de l'exécution du script à $(date)"

    # Vérifie si imagemagick, jq, fim et libcamera-jpeg sont installés
    for cmd in convert jq fim libcamera-jpeg; do
        if ! command -v $cmd &> /dev/null; then
            echo "$cmd n'est pas installé. Installation..."
            sudo apt-get install $cmd -y
        fi
    done

    # Vérifie si les polices Arial sont disponibles
    if ! fc-list | grep -qi "Arial"; then
        echo "Les polices Arial ne sont pas installées. Installation..."
        sudo apt-get install ttf-mscorefonts-installer -y
        sudo fc-cache -f -v
    fi

    # Prendre une photo
    libcamera-jpeg -o photo.jpg --rotation 180

    # Définir les variables
    API_KEY="${OPENAI_API_KEY}" # Utilise la variable d'environnement
    if [ -z "$API_KEY" ]; then
        echo "Erreur : La clé API n'est pas définie."
        exit 1
    fi
    MODEL="gpt-4o"
    IMAGE_PATH="photo.jpg"
    DATE=$(date +"%d %B %Y")
    TIME=$(date +"%Hh%M")

    # Encoder l'image en base64 et stocker dans une variable
    IMAGE_BASE64=$(base64 -w 0 $IMAGE_PATH)

    # Construire le contenu dynamique
    CONTENT="Il est ${TIME}, le ${DATE}. Tu es une caméra sur un Raspberry Pi dans mon salon, et tu observes de temps en temps ce qu’il s’y passe pour afficher un mot sur l’écran du Raspberry. Ton travail c’est d’observer la photo que je t’envoie et de trouver quelque chose d’agréable et positif à écrire sur l’écran. Tu formules uniquement une phrase courte, en français, positive, liée à ce que tu vois sur la photo, sympathique, family friendly pour donner de l’amour à ceux qui te lisent. Sois précis en décrivant une personne, une activité ou une action visible sur la photo."

    # Créer un fichier temporaire pour la requête JSON
    REQUEST_PAYLOAD=$(mktemp)
    cat <<EOF > $REQUEST_PAYLOAD
{
  "model": "$MODEL",
  "messages": [
    {
      "role": "system",
      "content": "$CONTENT"
    },
    {
      "role": "user",
      "content": [
        {
          "type": "image_url",
          "image_url": {
            "url": "data:image/jpeg;base64,${IMAGE_BASE64}"
          }
        }
      ]
    }
  ],
  "temperature": 1,
  "max_tokens": 256,
  "top_p": 1,
  "frequency_penalty": 0,
  "presence_penalty": 0
}
EOF

    # Faire la requête CURL et récupérer la réponse
    RESPONSE=$(curl -s https://api.openai.com/v1/chat/completions \
      -H "Content-Type: application/json" \
      -H "Authorization: Bearer $API_KEY" \
      -d @$REQUEST_PAYLOAD)

    # Supprimer le fichier temporaire
    rm $REQUEST_PAYLOAD

    # Afficher la réponse brute pour debug
    echo "Réponse brute: $RESPONSE"

    # Extraire le texte de la réponse en utilisant jq
    CAPTION=$(echo $RESPONSE | jq -r '.choices[0].message.content')
    echo "Texte extrait: $CAPTION"

    # Créer une image avec du padding et une police sympa
    convert -background black -fill white -font Arial -pointsize 72 \
    -gravity southwest -extent 1280x720 -size 1200x600 caption:"$CAPTION" \
    -bordercolor black -border 100x100 -gravity southwest -extent 1280x720+50+50 output.png

    # Afficher l'image en plein écran
    fim -a --quiet output.png

    echo "Fin de l'exécution du script à $(date)"
} | sudo tee -a $LOG_FILE > /dev/null