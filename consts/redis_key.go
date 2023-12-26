package consts

const (
	REDIS_KEY_REQUEST_SESSION  = "nyacaptcha:pending:%s"
	REDIS_KEY_RESOLVED_SESSION = "nyacaptcha:resolved:%s"
)

const (
	REDIS_KEY_IP_COOLDOWN = "nyacaptcha:ipcd:%s:%s"
)
