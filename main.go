package main

import (
	"fmt"
	"go.uber.org/zap"
	"log"
	"nya-captcha/config"
	g "nya-captcha/global"
	"nya-captcha/inits"
	"nya-captcha/injects"
)

func main() {
	log.Println(fmt.Sprintf("正在启动 NyaCap Server Mini (%s)...", injects.VERSION))

	// 初始化配置
	if err := inits.Config(); err != nil {
		log.Fatalln("配置加载失败:", err)
	}

	// 初始化 logger
	if err := inits.Logger(); err != nil {
		log.Fatalln("Logger 初始化失败:", err)
	}

	g.Logger.Info("Logger 初始化完成，执行切换。")

	// 初始化 Redis
	if err := inits.Redis(); err != nil {
		g.Logger.Fatal("Redis 初始化失败", zap.Error(err))
	}

	// 初始化验证码核心
	if err := inits.Captcha(); err != nil {
		g.Logger.Fatal("验证码核心加载失败", zap.Error(err))
	}

	// 初始化 HTTP Server
	e := inits.WebEngine()

	g.Logger.Info("初始化完成")

	// Start
	g.Logger.Info("Service starting...")
	if err := e.Start(config.Config.System.Listen); err != nil {
		g.Logger.Fatal("服务启动失败", zap.Error(err))
	}
}
