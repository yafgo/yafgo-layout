package global

import (
	"yafgo/yafgo-layout/pkg/notify"
)

func SetupNotify() {
	notify.SetupFeishu(AppName(), AppEnv(), FeishuRobot())
}

func FeishuRobot() string {
	return Ycfg.GetString("feishu.default")
}
