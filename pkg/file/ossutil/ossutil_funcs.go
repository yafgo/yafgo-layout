package ossutil

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

func (p *OssUtil) Client() *oss.Client {
	return p.client
}
