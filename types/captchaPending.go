package types

import "github.com/wenlng/go-captcha/captcha"

type CaptchaPending struct {
	Origin    string                  `json:"origin"`
	IP        string                  `json:"ip"`
	UserAgent string                  `json:"ua"`
	Dots      map[int]captcha.CharDot `json:"dots"`
}
