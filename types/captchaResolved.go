package types

type CaptchaResolved struct {
	Success     bool   `json:"success"`
	IP          string `json:"ip"`
	Site        string `json:"site"`
	ChallengeTS string `json:"challenge_ts"` // (ISO format yyyy-MM-dd'T'HH:mm:ssZZ)
}
