package types

type CaptchaResolved struct {
	Success     bool     `json:"success"`
	IP          *string  `json:"ip,omitempty"`
	Origin      *string  `json:"origin,omitempty"`
	ChallengeTS *string  `json:"challenge_ts,omitempty"` // (ISO format)
	Hostname    *string  `json:"hostname,omitempty"`
	ErrorCodes  []string `json:"error-codes,omitempty"`
}
