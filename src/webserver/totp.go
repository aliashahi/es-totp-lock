package webserver

import (
	"encoding/base32"
	googauth "es-project/src/googauth2"
	"fmt"
	"math/rand"
)

func (u *User) Validate(passcode string) (bool, error) {
	otpc := &googauth.OTPConfig{
		Secret:      u.Secret,
		WindowSize:  3,
		HotpCounter: 0,
		UTC:         false,
	}

	val, err := otpc.Authenticate(passcode)
	if err != nil {
		return false, err
	}

	return val, nil
}

func GetUserByPasscode(passcode string) (*User, error) {
	for _, u := range users {
		if ok, err := u.Validate(passcode); err != nil {
			return nil, err
		} else if ok {
			return u, nil
		}
	}

	return nil, fmt.Errorf("wrong code")
}

func createSecret() string {
	var s string
	const chars = "qwertyuioplkjhgfdsamnbvcxzQWERTYUIOPASDFGHJKLZXCVBNM"
	for i := 0; i < 20; i++ {
		r := rand.Int31n(int32(len(chars)))
		s = fmt.Sprintf("%s%c", s, chars[r])
	}

	return base32.StdEncoding.EncodeToString([]byte(s))
}

func (u *User) IsAdmin() bool {
	return u.Username == "admin"
}
