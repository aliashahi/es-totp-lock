package webserver

import (
	"encoding/base32"
	googauth "es-project/src/googauth2"
	"fmt"
	"math/rand"
)

func (u *User) Validate(passcode string) (bool, error) {
	for {
		if len(passcode) < 6 {
			passcode = "0" + passcode
		} else {
			break
		}
	}
	otpc := &googauth.OTPConfig{
		Secret:      u.Secret,
		WindowSize:  3,
		HotpCounter: 0,
		UTC:         true,
	}

	val, err := otpc.Authenticate(passcode)
	if err != nil {
		return false, err
	}

	return val, nil
}

func convertToString(passcode []byte) string {
	v := ""
	for _, c := range passcode {
		v += fmt.Sprint(int(c) - 48)
	}
	return v
}

func GetUserByPasscode(passcode []byte) (*User, error) {
	code := convertToString(passcode)

	for _, u := range users {
		if ok, err := u.Validate(code); err != nil {
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
