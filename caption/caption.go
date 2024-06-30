package caption

import (
	"fmt"
	"time"

	"github.com/tidwall/gjson"
)

type Payload struct {
    Model     string        `json:"model"`
    Messages  []Message     `json:"messages"`
    Temperature float64     `json:"temperature"`
    MaxTokens int           `json:"max_tokens"`
    TopP      float64       `json:"top_p"`
    FrequencyPenalty float64 `json:"frequency_penalty"`
    PresencePenalty float64  `json:"presence_penalty"`
}

type Message struct {
    Role    string      `json:"role"`
    Content interface{} `json:"content"`
}

func CreateContent() string {
    date := time.Now().Format("02 January 2006")
    time := time.Now().Format("15h04")
    return fmt.Sprintf("Il est %s, le %s. Tu es une caméra sur un Raspberry Pi dans mon salon, et tu observes de temps en temps ce qu’il s’y passe pour afficher un mot sur l’écran du Raspberry. Ton travail c’est d’observer la photo que je t’envoie et de trouver quelque chose d’agréable et positif à écrire sur l’écran. Tu formules uniquement une phrase courte, en français, positive, liée à ce que tu vois sur la photo, sympathique, family friendly pour donner de l’amour à ceux qui te lisent. Sois précis en décrivant une personne, une activité ou une action visible sur la photo.", time, date)
}

func CreatePayload(content string, imageBase64 string, model string) Payload {
    return Payload{
        Model: model,
        Messages: []Message{
            {
                Role: "system",
                Content: content,
            },
            {
                Role: "user",
                Content: []map[string]interface{}{
                    {
                        "type": "image_url",
                        "image_url": map[string]string{
                            "url": "data:image/jpeg;base64," + imageBase64,
                        },
                    },
                },
            },
        },
        Temperature: 1,
        MaxTokens: 256,
        TopP: 1,
        FrequencyPenalty: 0,
        PresencePenalty: 0,
    }
}

func ExtractCaption(responseBody []byte) string {
    return gjson.Get(string(responseBody), "choices.0.message.content").String()
}