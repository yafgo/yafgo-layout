package ycfg

import (
	"strings"

	"github.com/gookit/color"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

/*
# config.yaml nacos 配置示例
nacos:
  dir_log: "./storage/tmp/nacos/log"
  dir_cache: "./storage/tmp/nacos/cache"
  config:
    host: "127.0.0.1"
    port: 8848
    context_path: "/nacos"
    namespace: "158dde6c-66e7-4107-a0bb-f879d02c767e"
    log_level: "info"
    data_id: "yafgo.yaml"
    group: "DEFAULT_GROUP"
*/

// setupNacos 启用 Nacos 支持
func (p *Config) setupNacos() (err error) {
	nacosCfgKey := "nacos.config"
	if !p.IsSet("nacos") || !p.IsSet(nacosCfgKey) {
		color.Warnln("nacos配置不存在, 不启用nacos配置功能")
		return
	}

	// nacos client相关目录配置
	p.SetDefault("nacos", map[string]any{
		"dir_log":   "./storage/tmp/nacos/log",
		"dir_cache": "./storage/tmp/nacos/cache",
	})

	// nacos 配置中心参数
	p.SetDefault(nacosCfgKey, map[string]any{
		"host":         "127.0.0.1", //
		"port":         8848,        //
		"context_path": "/nacos",    //
		"namespace":    "",          //
		"log_level":    "info",      // debug,info,warn,error, default value is info
		"data_id":      "",          //
		"group":        "",          //
	})
	nacosCfg := p.Sub(nacosCfgKey)
	dataId := nacosCfg.GetString("data_id")
	dataGroup := nacosCfg.GetString("group")
	namespace := nacosCfg.GetString("namespace")
	logLevel := nacosCfg.GetString("log_level")

	// 创建 ServerConfig
	var sc = []constant.ServerConfig{
		*constant.NewServerConfig(nacosCfg.GetString("host"), nacosCfg.GetUint64("port"), constant.WithContextPath(nacosCfg.GetString("context_path"))),
	}
	// 创建 ClientConfig
	var cc = constant.ClientConfig{
		NamespaceId:         namespace, // 如果需要支持多namespace，我们可以创建多个client,它们有不同的NamespaceId。当namespace是public时，此处填空字符串。
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              p.GetString("nacos.dir_log"),
		CacheDir:            p.GetString("nacos.dir_cache"),
		LogLevel:            logLevel,
	}
	// 创建动态配置客户端
	client, err := clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  &cc,
			ServerConfigs: sc,
		},
	)
	if err != nil {
		color.Errorf("初始化nacos失败: %+v\n", err)
		return
	}

	// 获取配置
	content, err := client.GetConfig(vo.ConfigParam{
		DataId: dataId,
		Group:  dataGroup,
	})
	if err != nil {
		color.Errorf("GetConfig err:%+v\n", err)
		return
	}
	if _err := p.MergeConfig(strings.NewReader(content)); _err != nil {
		color.Errorf("读取nacos配置出错: %+v\n", _err)
	}

	// Listen config change,key=dataId+group+namespaceId.
	err = client.ListenConfig(vo.ConfigParam{
		DataId: dataId,
		Group:  dataGroup,
		OnChange: func(namespace, group, dataId, data string) {
			color.Successf("nacos config changed, group:%s, dataId: %s", group, dataId)
			r := strings.NewReader(data)
			if _err := p.MergeConfig(r); _err != nil {
				color.Errorf("读取nacos配置出错: %+v\n", _err)
				return
			}
			p.EmitChange()
		},
	})
	if err != nil {
		color.Errorf("监听nacos配置出错: %+v\n", err)
	}
	color.Successf("MergeInConfig: %+v\n", nacosCfg.AllSettings())

	return err
}
