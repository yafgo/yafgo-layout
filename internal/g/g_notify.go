package g

func SetupNotify() {
	// notify.SetupFeishu(AppName(), AppEnv(), FeishuRobot())
}

func FeishuRobot() string {
	return Cfg().GetString("feishu.default")
}
