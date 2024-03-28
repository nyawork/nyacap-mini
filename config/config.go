package config

import (
	"nya-captcha/types"
	"time"
)

type cfgType struct {
	Sites   []types.SiteInfo `yaml:"sites"`
	Captcha struct {
		PendingValidFor time.Duration `yaml:"pending_valid_for"` // 申请后的有效时间
		SubmitValidFor  time.Duration `yaml:"submit_valid_for"`  // 提交后的会话有效时间
		Characters      []string      `yaml:"characters"`        // 有效字符
	} `yaml:"captcha"`
	Security struct {
		IPBanPeriod            time.Duration `yaml:"ip_ban_period"`      // 请求不匹配时 ban IP 的时间
		CaptchaRequestCooldown time.Duration `yaml:"captcha_request_cd"` // 重新请求验证码的间隔时间
		CaptchaSubmitCooldown  time.Duration `yaml:"captcha_submit_cd"`  // 提交后再次提交的冷却时间
	} `yaml:"security"`
	System struct {
		Debug bool   `yaml:"debug"`
		Redis string `yaml:"redis"`
	} `yaml:"system"`
}

var Config cfgType
