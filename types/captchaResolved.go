package types

type CaptchaResolved struct {
	// 模仿 recaptcha 的格式，方便后端迁移使用
	Success     bool     `json:"success"`
	ChallengeTS string   `json:"challenge_ts"`          // (ISO format yyyy-MM-dd'T'HH:mm:ssZZ)
	Hostname    string   `json:"hostname"`              // 从 Origin 转换
	ErrorCodes  []string `json:"error-codes,omitempty"` // 保留这个字段，但不会提供数据

	// 扩展
	IP   string `json:"ip"`
	Site string `json:"site"`
}
