package config

import (
	"nyacap-mini/types"
	"time"
)

type cfgType struct {
	Sites   []types.SiteInfo `yaml:"sites"`
	Captcha struct {
		PendingValidFor time.Duration `yaml:"pending_valid_for"` // 申请后的有效时间
		SubmitValidFor  time.Duration `yaml:"submit_valid_for"`  // 提交后的会话有效时间
		Characters      []string      `yaml:"characters"`        // 有效字符
		Padding         int64         `yaml:"padding"`           // 允许的误差范围
		CheckTextLen    struct {
			Min int `yaml:"min"`
			Max int `yaml:"max"`
		} `yaml:"check_text_len"` // 检查字符数量
	} `yaml:"captcha"`
	Security struct {
		IPBanPeriod            time.Duration `yaml:"ip_ban_period"`      // 请求不匹配时 ban IP 的时间
		CaptchaRequestCooldown time.Duration `yaml:"captcha_request_cd"` // 重新请求验证码的间隔时间
		CaptchaSubmitCooldown  time.Duration `yaml:"captcha_submit_cd"`  // 提交后再次提交的冷却时间
	} `yaml:"security"`
	System struct {
		Debug  bool   `yaml:"debug"`
		Redis  string `yaml:"redis"`
		Listen string `yaml:"listen"`
	} `yaml:"system"`
}

var Config cfgType
