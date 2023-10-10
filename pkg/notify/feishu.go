package notify

import (
	"context"
	"encoding/json"
	"strings"
	"time"
	"yafgo/yafgo-layout/pkg/sys/ycfg"
	"yafgo/yafgo-layout/pkg/sys/ylog"

	"github.com/imroc/req/v3"
)

// H is a shortcut for map[string]interface{}
type H map[string]any

type FeishuRobot struct {
	logger   *ylog.Logger
	prefix   string // 消息前缀
	robotUrl string // 默认robotUrl

	// 消息队列
	msgQueue chan *feishuMsg
	msgQSize int

	// 执行状态
	running bool
}

// NewFeishu
//
//	appEnv: 环境, dev, prod 等
//	robotUrl: 默认robot地址
func NewFeishu(logger *ylog.Logger, cfg *ycfg.Config) *FeishuRobot {
	fsRobot := new(FeishuRobot)
	fsRobot.logger = logger
	fsRobot.robotUrl = cfg.GetString("feishu.robot.default")
	fsRobot.prefix = cfg.GetString("feishu.prefix")

	if fsRobot.msgQSize <= 0 {
		fsRobot.msgQSize = 10
	}
	fsRobot.msgQueue = make(chan *feishuMsg, fsRobot.msgQSize)
	return fsRobot
}

func (fr *FeishuRobot) New() *feishuMsg {
	msg := new(feishuMsg)
	msg.feishuRobot = fr
	msg.robotUrl = fr.robotUrl
	return msg
}

func (fr *FeishuRobot) startFeishuMsgConsumer() {
	if fr.running {
		return
	}
	fr.running = true
	go func() {
		fr.logger.Infof(context.Background(), "启动飞书队列")
		for msg := range fr.msgQueue {
			if msg == nil {
				continue
			}

			go fr.doSendFeishuMsg(msg)
			time.Sleep(time.Millisecond * 100)
		}
	}()
}

func (fr *FeishuRobot) send(msg *feishuMsg) {
	if !fr.running {
		fr.startFeishuMsgConsumer()
	}
	fr.msgQueue <- msg
}

// send 发送
func (fr *FeishuRobot) doSendFeishuMsg(msg *feishuMsg) (err error) {
	ctx := context.Background()
	defer func() {
		if rec := recover(); rec != nil {
			fr.logger.Infof(ctx, "飞书消息发送错误defer: %+v", rec)
		}
	}()

	if msg == nil {
		return
	}

	if msg.robotUrl == "" {
		fr.logger.Errorf(ctx, "飞书机器人 webhook URL 未配置")
		return
	}

	if !strings.HasPrefix(msg.robotUrl, "http") {
		msg.robotUrl = "https://open.feishu.cn/open-apis/bot/v2/hook/" + msg.robotUrl
	}

	// 执行发送
	jsonReq, err := json.Marshal(msg.body)
	if err != nil {
		fr.logger.Infof(ctx, "飞书消息发送错误1: %+v", err)
		return
	}

	client := req.C().SetTimeout(10 * time.Second)
	resp, err := client.R().
		SetBody(string(jsonReq)).
		SetHeader("Content-Type", "application/json").
		Post(msg.robotUrl)
	if err != nil {
		fr.logger.Infof(ctx, "飞书消息发送错误2: %+v", err)
		return
	}
	if resp.IsErrorState() {
		fr.logger.Infof(ctx, "飞书消息发送失败, 响应:%s, 请求:%s", resp.String(), jsonReq)
	}
	return
}

func (fr *FeishuRobot) Stats() map[string]int {
	queued := len(fr.msgQueue)
	return map[string]int{
		"容量": fr.msgQSize,
		"队列": queued,
		"空闲": fr.msgQSize - queued,
	}
}
