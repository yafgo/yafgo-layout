package g

import (
	"yafgo/yafgo-layout/pkg/notify"
)

func SetupNotify() {
	notify.SetupFeishu(AppName(), AppEnv(), FeishuRobot())
}

func FeishuRobot() string {
	return Cfg.GetString("feishu.default")
}
