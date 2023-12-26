package inits

import (
	"github.com/wenlng/go-captcha/captcha"
	"nya-captcha/config"
	"nya-captcha/global"
	"os"
	"path/filepath"
)

func Captcha() error {
	var err error

	global.Captcha = captcha.NewCaptcha()

	// 初始化一些设定
	global.Captcha.SetRangCheckTextLen(captcha.RangeVal{
		Max: 5,
		Min: 3,
	})
	global.Captcha.SetTextShadow(true)
	global.Captcha.SetTextShadowPoint(captcha.Point{
		X: 3,
		Y: 3,
	})
	global.Captcha.SetTextShadowColor("#ffffff")
	global.Captcha.SetImageFontAlpha(1)
	global.Captcha.SetTextRangFontColors([]string{
		"#1e293b",
		"#1f2937",
		"#27272a",
		"#262626",
		"#292524",
		"#991b1b",
		"#9a3412",
		"#92400e",
		"#854d0e",
		"#3f6212",
		"#166534",
		"#065f46",
		"#115e59",
		"#155e75",
		"#075985",
		"#1e40af",
		"#3730a3",
		"#5b21b6",
		"#6b21a8",
		"#86198f",
		"#9d174d",
		"#9f1239",
	})

	//// 有效字符
	err = global.Captcha.SetRangChars(config.Config.Captcha.Characters)
	if err != nil {
		return err
	}

	// 获得当前工作目录，用于读取文件类配置
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	// 检查数据目录
	dataDirPath := filepath.Join(cwd, "data")
	_, err = os.Stat(dataDirPath)
	if err != nil {
		if os.IsNotExist(err) {
			// 数据目录不存在，使用默认配置信息
			return nil
		} else {
			// 无法获得数据目录的信息，返回错误
			return err
		}
	}

	// 加载字体文件
	fontFiles, err := listFileRecursive(filepath.Join(dataDirPath, "font"))
	if err != nil {
		return err
	}
	if len(fontFiles) > 0 {
		global.Captcha.SetFont(fontFiles)
	}

	// 加载背景图片
	backgroundImages, err := listFileRecursive(filepath.Join(dataDirPath, "background"))
	if err != nil {
		return err
	}
	if len(backgroundImages) > 0 {
		global.Captcha.SetBackground(backgroundImages)
	}

	// 加载小图
	thumbnailImages, err := listFileRecursive(filepath.Join(dataDirPath, "thumbnail"))
	if err != nil {
		return err
	}
	if len(thumbnailImages) > 0 {
		global.Captcha.SetThumbBackground(thumbnailImages)
	}

	// 初始化完成
	return nil

}

func listFileRecursive(root string) ([]string, error) {
	var fileNames []string
	contents, err := os.ReadDir(root)
	if err != nil {
		return nil, err
	}

	for _, content := range contents {
		if content.IsDir() {
			subContents, err := listFileRecursive(filepath.Join(root, content.Name()))
			if err != nil {
				return nil, err
			} else {
				fileNames = append(fileNames, subContents...)
			}
		} else {
			fileNames = append(fileNames, filepath.Join(root, content.Name()))
		}
	}

	return fileNames, nil

}
