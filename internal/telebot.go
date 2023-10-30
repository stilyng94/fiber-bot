package internal

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"slices"
	"strings"

	"github.com/fiber-bot/config"
)

type DecodeTelegramHashPayload struct {
	TelegramID int64  `json:"id" validate:"required"`
	FirstName  string `json:"first_name,omitempty"`
	LastName   string `json:"last_name,omitempty"`
	Username   string `json:"username,omitempty"`
	QueryID    string `json:"query_id" validate:"required"`
}

func DecodeTelegramHash(token string, envConfig config.EnvConfig) (DecodeTelegramHashPayload, error) {
	parsedToken, err := url.ParseQuery(token)
	if err != nil {
		return DecodeTelegramHashPayload{}, err
	}

	dataToCheck := []string{}
	for k, v := range parsedToken {
		if k == "hash" {
			continue
		}

		dataToCheck = append(dataToCheck, fmt.Sprintf("%s=%s", k, v[0]))
	}

	slices.Sort(dataToCheck)

	secretKey := hmac.New(sha256.New, []byte("WebAppData"))
	secretKey.Write([]byte(envConfig.TelegramBotToken))

	hash := hmac.New(sha256.New, secretKey.Sum(nil))
	hash.Write([]byte(strings.Join(dataToCheck, "\n")))

	calculatedHash := hex.EncodeToString(hash.Sum(nil))

	if calculatedHash != parsedToken.Get("hash") {
		return DecodeTelegramHashPayload{}, errors.New("invalid token")
	}

	var payload DecodeTelegramHashPayload
	err = json.Unmarshal([]byte(parsedToken.Get("user")), &payload)
	if err != nil {
		return DecodeTelegramHashPayload{}, errors.New("cannot parse user")
	}
	if parsedToken.Has("is_bot") && parsedToken.Get("is_bot") == "true" {
		return DecodeTelegramHashPayload{}, errors.New("bot's not allowed")
	}
	payload.QueryID = parsedToken.Get("query_id")
	return payload, nil
}
