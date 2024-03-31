package security

import (
	"context"
	"fmt"
	"nya-captcha/consts"
	g "nya-captcha/global"
	"time"
)

func CooldownIP(ip string, pool string, cd time.Duration) {
	g.Redis.Set(
		context.Background(),
		fmt.Sprintf(consts.REDIS_KEY_IP_COOLDOWN, pool, ip),
		nil,
		cd,
	)
}
