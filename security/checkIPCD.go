package security

import (
	"context"
	"fmt"
	"nya-captcha/consts"
	g "nya-captcha/global"
)

func CheckIPCD(ip string, pool string) (bool, error) {
	ipKey := fmt.Sprintf(consts.REDIS_KEY_IP_COOLDOWN, pool, ip)
	exist, err := g.Redis.Exists(context.Background(), ipKey).Result()
	return exist > 0, err
}
