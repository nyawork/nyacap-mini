package security

import (
	"context"
	"fmt"
	"nya-captcha/consts"
	"nya-captcha/global"
	"time"
)

func CooldownIP(ip string, pool string, cd time.Duration) {
	global.Redis.Set(
		context.Background(),
		fmt.Sprintf(consts.REDIS_KEY_IP_COOLDOWN, pool, ip),
		nil,
		cd,
	)
}
