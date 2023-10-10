package notify

import "fmt"

type feishuMsg struct {
	feishuRobot *FeishuRobot

	// robot地址
	robotUrl string
	// 消息内容
	body any

	// 其他字段
	atAll bool // 是否 @所有人
}

// WithRobot 本次消息发给指定的飞书 robot
func (p *feishuMsg) WithRobot(robotUrl string) *feishuMsg {
	p.robotUrl = robotUrl
	return p
}

// AtAll 本次消息@所有人
func (p *feishuMsg) AtAll() *feishuMsg {
	p.atAll = true
	return p
}

// send 发送
func (p *feishuMsg) send() {
	p.feishuRobot.send(p)
}

// SendText 发送文本消息
//
//	SendText("测试发送文本消息")
//	SendText("测试发送文本消息: %s, %d", "张三", 18)
func (p *feishuMsg) SendText(text string, a ...any) {
	if len(a) > 0 {
		text = fmt.Sprintf(text, a...)
	}
	if p.feishuRobot.prefix != "" {
		text = p.feishuRobot.prefix + " " + text
	}
	if p.atAll {
		text += ` <at user_id="all">所有人</at>`
	}

	// 发送数据
	p.body = H{
		"msg_type": "text",
		"content": H{
			"text": text,
		},
	}
	p.send()
}

// SendPost 发送富文本消息
//
//	SendPost("测试标题", "测试发送文本消息")
//	SendPost("测试标题", "测试发送文本消息: %s, %d", "张三", 18)
func (p *feishuMsg) SendPost(title string, text string, a ...any) {
	if len(a) > 0 {
		text = fmt.Sprintf(text, a...)
	}
	if p.feishuRobot.prefix != "" {
		title = p.feishuRobot.prefix + " " + title
	}
	var content = []H{
		{"tag": "text", "text": text},
	}
	if p.atAll {
		content = append(content, H{"tag": "at", "user_id": "all"})
	}
	var contents = [][]H{content}

	// 发送数据
	p.body = H{
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
	p.send()
}
