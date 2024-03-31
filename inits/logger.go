package inits

import (
	"fmt"
	"go.uber.org/zap"
	"nya-captcha/config"
	g "nya-captcha/global"
)

func Logger() error {
	var err error

	// Prepare logger
	if config.Config.System.Debug {
		g.Logger, err = zap.NewDevelopment()
	} else {
		g.Logger, err = zap.NewProduction()
	}
	if err != nil {
		return fmt.Errorf("logger 初始化失败: %w", err)
	}

	return nil
}
