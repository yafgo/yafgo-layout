package notify

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"
	"yafgo/yafgo-layout/pkg/sys/ylog"

	"github.com/imroc/req/v3"
)

// H is a shortcut for map[string]interface{}
type H map[string]any

type feishuUtil struct {
	appEnv   string
	robotUrl string

	// 其他字段
	atAll bool // 是否 @所有人
}

type feishuMsg struct {
	robotUrl string
	body     any
}

var feishuQueueSize = 3000
var feishuMsgQueue = make(chan *feishuMsg, feishuQueueSize)

func Stats() map[string]int {
	queued := len(feishuMsgQueue)
	return map[string]int{
		"容量": feishuQueueSize,
		"队列": queued,
		"空闲": feishuQueueSize - queued,
	}
}

func startFeishuMsgConsumer() {
	go func() {
		ylog.Infof(context.Background(), "启动飞书队列")
		for msg := range feishuMsgQueue {
			if msg == nil {
				continue
			}

			go sendFeishuMsg(msg)
			time.Sleep(time.Millisecond * 200)
		}
	}()
}

// send 发送
func sendFeishuMsg(msg *feishuMsg) (err error) {
	ctx := context.Background()
	defer func() {
		if rec := recover(); rec != nil {
			ylog.Infof(ctx, "飞书消息发送错误defer: %+v", rec)
		}
	}()

	if msg == nil {
		return
	}

	if msg.robotUrl == "" {
		ylog.Errorf(ctx, "飞书机器人 webhook URL 未配置")
		return
	}

	if !strings.HasPrefix(msg.robotUrl, "http") {
		msg.robotUrl = "https://open.feishu.cn/open-apis/bot/v2/hook/" + msg.robotUrl
	}

	// 执行发送
	jsonReq, err := json.Marshal(msg.body)
	if err != nil {
		ylog.Infof(ctx, "飞书消息发送错误1: %+v", err)
		return
	}

	client := req.C().SetTimeout(10 * time.Second)
	resp, err := client.R().
		SetBody(string(jsonReq)).
		SetHeader("Content-Type", "application/json").
		Post(msg.robotUrl)
	if err != nil {
		ylog.Infof(ctx, "飞书消息发送错误2: %+v", err)
		return
	}
	if resp.IsErrorState() {
		ylog.Infof(ctx, "飞书消息发送失败, 响应:%s, 请求:%s", resp.String(), jsonReq)
	}
	return
}

var (
	defaultAppEnv   string
	defaultRobotUrl string
	_sysName        string
)

// SetupFeishu 初始化飞书配置
//
//	robotUrl: 飞书机器人 webhook v2 地址, 也可以省略 url 前缀, 只写最后一段: "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
func SetupFeishu(sysName, appEnv, robotUrl string) {
	_sysName = sysName
	defaultAppEnv = appEnv
	defaultRobotUrl = robotUrl
	startFeishuMsgConsumer()
}

func Feishu() (fsUtil *feishuUtil) {
	fsUtil = new(feishuUtil)
	fsUtil.appEnv = defaultAppEnv
	fsUtil.robotUrl = defaultRobotUrl
	return
}

// WithRobot 本次消息发给指定的飞书 robot
func (p *feishuUtil) WithRobot(robotUrl string) *feishuUtil {
	p.robotUrl = robotUrl
	return p
}

// AtAll 本次消息@所有人
func (p *feishuUtil) AtAll() *feishuUtil {
	p.atAll = true
	return p
}

// send 发送
func (p *feishuUtil) send(body any) (err error) {
	feishuMsgQueue <- &feishuMsg{p.robotUrl, body}
	return
}

// SendText 发送文本消息
//
//	SendText("测试发送文本消息")
//	SendText("测试发送文本消息: %s, %d", "张三", 18)
func (p *feishuUtil) SendText(text string, a ...any) error {
	if len(a) > 0 {
		text = fmt.Sprintf(text, a...)
	}
	text = fmt.Sprintf("【%s】%s", p.appEnv, text)
	if p.atAll {
		text += ` <at user_id="all">所有人</at>`
	}

	// 发送数据
	msgData := H{
		"msg_type": "text",
		"content": H{
			"text": text,
		},
	}
	return p.send(msgData)
}

// SendPost 发送富文本消息
//
//	SendPost("测试标题", "测试发送文本消息")
//	SendPost("测试标题", "测试发送文本消息: %s, %d", "张三", 18)
func (p *feishuUtil) SendPost(title string, text string, a ...any) error {
	if len(a) > 0 {
		text = fmt.Sprintf(text, a...)
	}
	title = fmt.Sprintf("【%s-%s】%s", _sysName, p.appEnv, title)
	var content = []H{
		{"tag": "text", "text": text},
	}
	if p.atAll {
		content = append(content, H{"tag": "at", "user_id": "all"})
	}
	var contents = [][]H{content}

	// 发送数据
	msgData := H{
		"msg_type": "post",
		"content": H{
			"post": H{
				"zh_cn": H{
					"title":   title,
					"content": contents,
				},
			},
		},
	}
	return p.send(msgData)
}
