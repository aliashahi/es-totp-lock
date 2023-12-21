package webserver

import (
	"encoding/base32"
	"time"

	"github.com/gorilla/websocket"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
)

func init() {
	updateCodes()
	go func() {
		for {
			updateCodes()
			time.Sleep(15 * time.Second)
		}
	}()
}

func updateCodes() {
	for _, u := range users {
		u.Code = generatePassCode(u.Secret)
		for _, client := range clientChannels {
			if client.UserID == u.ID {
				client.Conn.WriteMessage(websocket.TextMessage, []byte(u.Code))
			}
		}
	}
}

func generatePassCode(raw_secret string) string {
	secret := base32.StdEncoding.EncodeToString([]byte(raw_secret))
	passcode, err := totp.GenerateCodeCustom(secret, time.Now(), totp.ValidateOpts{
		Period:    30,
		Skew:      1,
		Digits:    otp.DigitsSix,
		Algorithm: otp.AlgorithmSHA512,
	})

	if err != nil {
		panic(err)
	}

	return passcode
}
