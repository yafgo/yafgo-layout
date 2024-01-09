// Package ossutil 封装阿里云oss相关操作
package ossutil

import (
	"fmt"
	"yafgo/yafgo-layout/pkg/sys/ycfg"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

type OssUtil struct {
	client *oss.Client
	cfg    *ycfg.Config

	endpoint        string
	accessKeyId     string
	accessKeySecret string
	bucket          string
}

// NewOssUtil 自定义封装的oss操作类
func NewOssUtil(cfg *ycfg.Config) (p *OssUtil) {
	p = new(OssUtil)
	p.cfg = cfg

	ossCfg := cfg.Sub("alioss")
	if ossCfg == nil {
		panic(fmt.Errorf("aliOss config 不存在"))
	}

	/*
	  # 阿里云Oss配置
	  alioss:
	    accessKeyId: ""
	    accessKeySecret: ""
	    bucketName: ""
	    endpoint: oss-cn-shanghai.aliyuncs.com
	*/
	p.endpoint = ossCfg.GetString("endpoint")
	p.accessKeyId = ossCfg.GetString("accessKeyId")
	p.accessKeySecret = ossCfg.GetString("accessKeySecret")
	p.bucket = ossCfg.GetString("bucketName")

	client, err := oss.New(p.endpoint, p.accessKeyId, p.accessKeySecret)
	if err != nil {
		panic(fmt.Errorf("aliOss Setup: %+v", err))
	}
	p.client = client
	return p
}
